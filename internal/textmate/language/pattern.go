// Copyright 2016 The lime Authors.
// Use of this source code is governed by a 2-clause
// BSD-style license that can be found in the LICENSE file.

package language

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/limetext/backend/log"
	"github.com/limetext/sublime/internal"
	"github.com/limetext/text"
	"github.com/quarnster/parser"
)

type (
	Pattern struct {
		internal.Named
		Include        string
		Match          internal.Regex
		Captures       internal.Captures
		Begin          internal.Regex
		BeginCaptures  internal.Captures
		End            internal.Regex
		EndCaptures    internal.Captures
		Patterns       []Pattern
		owner          *Language // needed for include directives
		cachedData     string
		cachedPat      *Pattern
		cachedPatterns []*Pattern // cached sub patterns
		cachedMatch    internal.MatchObject
		hits           int
		misses         int
	}

	RootPattern struct {
		Pattern
	}
)

var failed = make(map[string]bool)

func (r *RootPattern) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, &r.Patterns)
}

func (r *RootPattern) String() (ret string) {
	for i := range r.Patterns {
		ret += fmt.Sprintf("\t%s\n", r.Patterns[i])
	}
	return
}

// Finds the first sub pattern that has match for data after pos
func (p *Pattern) FirstMatch(data string, pos int) (pat *Pattern, ret internal.MatchObject) {
	startIdx := -1
	for i := 0; i < len(p.cachedPatterns); {
		ip, im := p.cachedPatterns[i].Cache(data, pos)
		// If it wasn't found now, it'll never be found,
		// so the pattern can be popped from the cache
		if im == nil {
			copy(p.cachedPatterns[i:], p.cachedPatterns[i+1:])
			p.cachedPatterns = p.cachedPatterns[:len(p.cachedPatterns)-1]
			continue
		}
		// ???: what is startIdx > im[0] for?
		if startIdx < 0 || startIdx > im[0] {
			startIdx, pat, ret = im[0], ip, im
			// This match is right at the start, we're not
			// going to find a better pattern than this,
			// so stop the search
			if im[0] == pos {
				break
			}
		}
		i++
	}
	return
}

func (p *Pattern) initCache() {
	if p.cachedPatterns != nil {
		return
	}
	p.cachedPatterns = make([]*Pattern, len(p.Patterns))
	for i := range p.cachedPatterns {
		p.cachedPatterns[i] = &p.Patterns[i]
	}
}

// Finds what does this pattern match also caches the match for next uses.
// Searches in order Match, Begin, Include, sub patterns.
func (p *Pattern) Cache(data string, pos int) (pat *Pattern, ret internal.MatchObject) {
	if p.cachedData == data {
		if p.cachedMatch == nil {
			return nil, nil
		}
		if p.cachedMatch[0] >= pos && p.cachedPat.cachedMatch != nil {
			p.hits++
			return p.cachedPat, p.cachedMatch
		}
	} else {
		p.cachedPatterns = nil
	}
	p.initCache()
	p.misses++

	if !p.Match.Empty() {
		pat, ret = p, p.Match.Find(data, pos)
	} else if !p.Begin.Empty() {
		pat, ret = p, p.Begin.Find(data, pos)
	} else if p.Include != "" {
		if z := p.Include[0]; z == '#' {
			key := p.Include[1:]
			if p2, ok := p.owner.Repository[key]; ok {
				pat, ret = p2.Cache(data, pos)
			} else {
				log.Fine("Not found in %s repository: %s", p.owner.Name, p.Include)
			}
		} else if z == '$' {
			// TODO(q): Implement tmLanguage $ include directives
			log.Warn("Unhandled include directive: %s", p.Include)
		} else if l, err := provider.GetLanguage(p.Include); err != nil {
			if !failed[p.Include] {
				log.Warn("Include directive %s failed: %s", p.Include, err)
			}
			failed[p.Include] = true
		} else {
			pat, ret = l.RootPattern.Cache(data, pos)
		}
	} else {
		pat, ret = p.FirstMatch(data, pos)
	}
	p.cachedData = data
	p.cachedMatch = ret
	p.cachedPat = pat

	return
}

// Creates a node for each capture also takse care of parent child relationship
func (p *Pattern) CreateCaptureNodes(data string, pos int, d parser.DataSource,
	mo internal.MatchObject, parent *parser.Node, capt internal.Captures) {
	// each couple of match objects will be a range
	ranges := make([]text.Region, len(mo)/2)
	// maps indexes from child to parent, the default int value is 0 so
	// by default each child parent will be the parent node
	parentIndex := make([]int, len(ranges))
	parents := make([]*parser.Node, len(parentIndex))
	for i := range ranges {
		// converting each couple of match objects to region
		ranges[i] = text.Region{A: mo[i*2], B: mo[i*2+1]}
		// the first 2 elements of parents node should be parent so we
		// could append childs easier later
		if i < 2 {
			parents[i] = parent
			continue
		}
		// if each pervious ranges covers current range we will store it
		// in parentIndex variable, because the node with wider range
		// should be the parent
		for j := i - 1; j >= 0; j-- {
			if ranges[j].Covers(ranges[i]) {
				parentIndex[i] = j
				break
			}
		}
	}

	for _, v := range capt {
		i := v.Key
		// ???: when do we set a range to -1?
		if i >= len(parents) || ranges[i].A == -1 {
			continue
		}
		// creating a node for each capture
		child := &parser.Node{Name: v.Name, Range: ranges[i], P: d}
		parents[i] = child
		if i == 0 {
			parent.Append(child)
			continue
		}
		var n *parser.Node
		// searches for child node parent then appends it
		for n == nil {
			i = parentIndex[i]
			n = parents[i]
		}
		n.Append(child)
	}
}

// Creates a root node for the pattern then creates a node for each capture and
// appends them as child of root node
func (p *Pattern) CreateNode(data string, pos int, d parser.DataSource, mo internal.MatchObject) (ret *parser.Node) {
	ret = &parser.Node{Name: p.Name, Range: text.Region{A: mo[0], B: mo[1]}, P: d}
	defer ret.UpdateRange()

	if !p.Match.Empty() {
		p.CreateCaptureNodes(data, pos, d, mo, ret, p.Captures)
	}
	if p.Begin.Empty() {
		return
	}
	if len(p.BeginCaptures) > 0 {
		p.CreateCaptureNodes(data, pos, d, mo, ret, p.BeginCaptures)
	} else {
		p.CreateCaptureNodes(data, pos, d, mo, ret, p.Captures)
	}

	if p.End.Empty() {
		return
	}
	var (
		found  = false
		i, end int
	)
	for i, end = ret.Range.B, len(data); i < len(data); {
		endmatch := p.End.Find(data, i)
		if endmatch == nil {
			if !found {
				// oops.. no end found at all, set it to the next line
				if e2 := strings.IndexRune(data[i:], '\n'); e2 != -1 {
					end = i + e2
				} else {
					end = len(data)
				}
				break
			}
			end = i
			break
		}
		end = endmatch[1]

		if len(p.cachedPatterns) > 0 {
			// Might be more recursive patterns to apply before the end is reached
			pattern2, match2 := p.FirstMatch(data, i)
			if match2 != nil && (match2[0] < endmatch[0] || (match2[0] == endmatch[0] && ret.Range.A == ret.Range.B)) {
				found = true
				r := pattern2.CreateNode(data, i, d, match2)
				ret.Append(r)
				i = r.Range.B
				continue
			}
		}
		if len(p.EndCaptures) > 0 {
			p.CreateCaptureNodes(data, i, d, endmatch, ret, p.EndCaptures)
		} else {
			p.CreateCaptureNodes(data, i, d, endmatch, ret, p.Captures)
		}
		break
	}
	ret.Range.B = end
	return
}

func (p *Pattern) copy(l *Language) *Pattern {
	ret := &Pattern{}
	ret.Named = p.Named
	ret.Include = p.Include
	ret.Match = *p.Match.Copy()
	if p.Captures != nil {
		ret.Captures = *p.Captures.Copy()
	}
	ret.Begin = *p.Begin.Copy()
	if p.BeginCaptures != nil {
		ret.BeginCaptures = *p.BeginCaptures.Copy()
	}
	ret.End = *p.End.Copy()
	if p.EndCaptures != nil {
		ret.EndCaptures = *p.EndCaptures.Copy()
	}
	ret.owner = l
	for _, pat := range p.Patterns {
		ret.Patterns = append(ret.Patterns, *pat.copy(l))
	}
	return ret
}

func (p *Pattern) tweak(l *Language) {
	p.owner = l
	p.Name = strings.TrimSpace(p.Name)
	for i := range p.Patterns {
		p.Patterns[i].tweak(l)
	}
}
