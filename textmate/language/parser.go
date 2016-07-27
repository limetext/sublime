// Copyright 2016 The lime Authors.
// Use of this source code is governed by a 2-clause
// BSD-style license that can be found in the LICENSE file.

package language

import (
	"strings"

	"github.com/limetext/backend/log"
	"github.com/limetext/text"
	"github.com/quarnster/parser"
)

// implements parser.Parser + parser.DataSource
type Parser struct {
	l    *Language
	data []rune
}

func NewParser(l *Language, data []rune) *Parser {
	return &Parser{l, data}
}

func (p *Parser) Data(a, b int) string {
	a = text.Clamp(0, len(p.data), a)
	b = text.Clamp(0, len(p.data), b)
	return string(p.data[a:b])
}

// fix node range positions based on given look up table
func (p *Parser) patch(lut []int, node *parser.Node) {
	node.Range.A = lut[node.Range.A]
	node.Range.B = lut[node.Range.B]
	for _, child := range node.Children {
		p.patch(lut, child)
	}
}

const maxiter = 10000

// Creates a directed acyclic graph from the data using the language
// finds the match pattern starting from each character in data then creates a
// node from the matched pattern and appends it to the root node
func (p *Parser) Parse() (*parser.Node, error) {
	sdata := string(p.data)
	rn := parser.Node{P: p, Name: p.l.ScopeName}
	defer rn.UpdateRange()
	defer func() {
		if r := recover(); r != nil {
			log.Error("Panic during parse: %v\n", r)
			log.Debug("%v", rn)
		}
	}()
	iter := maxiter
	for i := 0; i < len(sdata) && iter > 0; iter-- {
		pat, ret := p.l.RootPattern.Cache(sdata, i)
		if ret == nil {
			break
		}
		nl := strings.IndexAny(sdata[i:], "\n\r")
		// if the match object is not in this line we will move the
		// position to beginning of the next line then we will look for
		// a match again
		if nl += i; nl != i-1 && nl <= ret[0] {
			i = nl
			for i < len(sdata) && (sdata[i] == '\n' || sdata[i] == '\r') {
				i++
			}
			continue
		}
		n := pat.CreateNode(sdata, i, p, ret)
		rn.Append(n)
		i = n.Range.B
	}
	// handling unicode characters different length
	if len(sdata) != 0 {
		lut := make([]int, len(sdata)+1)
		j := 0
		for i := range sdata {
			lut[i] = j
			j++
		}
		lut[len(sdata)] = len(p.data)
		p.patch(lut, &rn)
	}
	if iter == 0 {
		panic("reached maximum number of iterations")
	}
	return &rn, nil
}
