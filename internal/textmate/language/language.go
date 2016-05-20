// Copyright 2013 The lime Authors.
// Use of this source code is governed by a 2-clause
// BSD-style license that can be found in the LICENSE file.

package language

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"strings"
	"sync"

	"github.com/limetext/loaders"
)

type (
	// For loading tmLanguage files
	// https://manual.macromates.com/en/language_grammars
	Language struct {
		UnpatchedLanguage
	}

	UnpatchedLanguage struct {
		FileTypes      []string
		FirstLineMatch string
		Name           string
		RootPattern    RootPattern `json:"patterns"`
		Repository     map[string]*Pattern
		ScopeName      string
	}

	Provider struct {
		sync.Mutex
		scope map[string]string
	}
)

func Load(filename string) (*Language, error) {
	d, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("Couldn't load file %s: %s", filename, err)
	}
	var l Language
	if err := loaders.LoadPlist(d, &l); err != nil {
		return nil, err
	}
	provider.Add(l.ScopeName, filename)
	return &l, nil
}

func (t *Provider) GetLanguage(id string) (*Language, error) {
	if l, err := t.LanguageFromScope(id); err != nil {
		return Load(id)
	} else {
		return l, err
	}
}

func (t *Provider) LanguageFromScope(id string) (*Language, error) {
	t.Lock()
	s, ok := t.scope[id]
	t.Unlock()
	if !ok {
		return nil, errors.New("Can't handle id " + id)
	} else {
		return Load(s)
	}
}

func (t *Provider) Add(scope, filename string) {
	t.Lock()
	defer t.Unlock()
	t.scope[scope] = filename
}

func (p Pattern) String() (ret string) {
	ret = fmt.Sprintf(`---------------------------------------
Name:    %s
Match:   %s
Begin:   %s
End:     %s
Include: %s
`, p.Name, p.Match, p.Begin, p.End, p.Include)
	ret += fmt.Sprintf("<Sub-Patterns>\n")
	for i := range p.Patterns {
		inner := fmt.Sprintf("%s", p.Patterns[i])
		ret += fmt.Sprintf("\t%s\n", strings.Replace(strings.Replace(inner, "\t", "\t\t", -1), "\n", "\n\t", -1))
	}
	ret += fmt.Sprintf("</Sub-Patterns>\n---------------------------------------")
	return
}

func (s *Language) String() string {
	return fmt.Sprintf("%s\n%s\n%s\n", s.ScopeName, s.RootPattern, s.Repository)
}

func (l *Language) tweak() {
	l.RootPattern.tweak(l)
	for k := range l.Repository {
		p := l.Repository[k]
		p.tweak(l)
		l.Repository[k] = p
	}
}

func (l *Language) UnmarshalJSON(data []byte) error {
	if err := json.Unmarshal(data, &l.UnpatchedLanguage); err != nil {
		return err
	}
	l.tweak()
	return nil
}

func (l *Language) Copy() *Language {
	ret := &Language{}
	ret.FileTypes = make([]string, len(l.FileTypes))
	copy(ret.FileTypes, l.FileTypes)
	ret.FirstLineMatch = l.FirstLineMatch
	ret.Name = l.Name
	ret.RootPattern.Pattern = *l.RootPattern.Pattern.copy(ret)
	ret.Repository = make(map[string]*Pattern)
	for key, pat := range l.Repository {
		ret.Repository[key] = pat.copy(ret)
	}
	ret.ScopeName = l.ScopeName
	return ret
}

var provider Provider

func init() {
	provider.scope = make(map[string]string)
}
