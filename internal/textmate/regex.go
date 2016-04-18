package internal

import (
	"fmt"
	"strings"

	"github.com/going/toolkit/log"
	"github.com/limetext/rubex"
)

type (
	Regex struct {
		re        *rubex.Regexp
		lastIndex int
		lastFound int
	}

	MatchObject []int
)

func (r Regex) Empty() bool {
	return r.re == nil
}

func (r Regex) String() string {
	if r.re == nil {
		return "nil"
	}
	return fmt.Sprintf("%s   // %d, %d", r.re.String(), r.lastIndex, r.lastFound)
}

func (r *Regex) UnmarshalJSON(data []byte) error {
	str := string(data[1 : len(data)-1])
	str = strings.Replace(str, "\\\\", "\\", -1)
	str = strings.Replace(str, "\\n", "\n", -1)
	str = strings.Replace(str, "\\t", "\t", -1)
	if re, err := rubex.Compile(str); err != nil {
		log.Warn("Couldn't compile language pattern %s: %s", str, err)
	} else {
		r.re = re
	}
	return nil
}

func (r *Regex) Find(data string, pos int) MatchObject {
	if r.lastIndex > pos {
		r.lastFound = 0
	}
	r.lastIndex = pos
	for r.lastFound < len(data) {
		ret := r.re.FindStringSubmatchIndex(data[r.lastFound:])
		if ret == nil {
			break
		} else if (ret[0] + r.lastFound) < pos {
			if ret[0] == 0 {
				r.lastFound++
			} else {
				r.lastFound += ret[0]
			}
			continue
		}
		mo := MatchObject(ret)
		mo.fix(r.lastFound)
		return mo
	}
	return nil
}

func (r *Regex) Copy() *Regex {
	ret := &Regex{}
	if r.re == nil {
		return ret
	}
	if re, err := rubex.Compile(fmt.Sprint(r.re)); err != nil {
		log.Warn("Error on copying regex: %s", err)
	} else {
		ret.re = re
	}
	ret.lastIndex = r.lastIndex
	ret.lastFound = r.lastFound
	return ret
}

func (m MatchObject) fix(add int) {
	for i := range m {
		if m[i] != -1 {
			m[i] += add
		}
	}
}
