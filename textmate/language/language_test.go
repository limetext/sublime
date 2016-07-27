// Copyright 2013 The lime Authors.
// Use of this source code is governed by a 2-clause
// BSD-style license that can be found in the LICENSE file.

package language

import (
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/limetext/util"
)

const gotmLang = "../../testdata/package/Go.tmLanguage"

func TestProviderLanguageFromScope(t *testing.T) {
	l, _ := Load(gotmLang)

	if _, err := provider.LanguageFromScope(l.ScopeName); err != nil {
		t.Errorf("Tried to load %s, but got an error: %v", l.ScopeName, err)
	}

	if _, err := provider.LanguageFromScope("MissingScope"); err == nil {
		t.Error("Tried to load MissingScope, expecting to get an error, but didn't")
	}
}

func TestProviderLanguageFromFile(t *testing.T) {
	if _, err := Load(gotmLang); err != nil {
		t.Errorf("Tried to load %s, but got an error: %v", gotmLang, err)
	}

	if _, err := Load("MissingFile"); err == nil {
		t.Error("Tried to load MissingFile, expecting to get an error, but didn't")
	}
}

func TestTmLanguage(t *testing.T) {
	files := []string{
		"testdata/Property List (XML).tmLanguage",
		"testdata/XML.plist",
		gotmLang,
	}
	for _, fn := range files {
		if _, err := Load(fn); err != nil {
			t.Fatal(err)
		}
	}

	tests := []struct {
		in  string
		out string
		syn string
	}{
		{
			"testdata/plist.tmlang",
			"testdata/plist.tmlang.res",
			"text.xml.plist",
		},
		{
			"testdata/Property List (XML).tmLanguage",
			"testdata/Property List (XML).tmLanguage.res",
			"text.xml.plist",
		},
		{
			"testdata/main.go",
			"testdata/main.go.res",
			"source.go",
		},
		{
			"testdata/go2.go",
			"testdata/go2.go.res",
			"source.go",
		},
		{
			"testdata/utf.go",
			"testdata/utf.go.res",
			"source.go",
		},
	}
	for _, test := range tests {
		var d0 string
		if d, err := ioutil.ReadFile(test.in); err != nil {
			t.Errorf("Couldn't load file %s: %s", test.in, err)
			continue
		} else {
			d0 = string(d)
		}

		if pr, err := getParser(test.syn, d0); err != nil {
			t.Error(err)
		} else if root, err := pr.Parse(); err != nil {
			t.Error(err)
		} else {
			str := fmt.Sprintf("%s", root)
			if d, err := ioutil.ReadFile(test.out); err != nil {
				if err := ioutil.WriteFile(test.out, []byte(str), 0644); err != nil {
					t.Error(err)
				}
			} else if diff := util.Diff(string(d), str); diff != "" {
				t.Errorf("%s:\n%s", test.in, diff)
			}
		}
	}
}

func BenchmarkLanguage(b *testing.B) {
	b.StopTimer()
	tst := []string{
		"testdata/utf.go",
		"testdata/main.go",
	}

	var d0 []string
	for _, s := range tst {
		if d, err := ioutil.ReadFile(s); err != nil {
			b.Errorf("Couldn't load file %s: %s", s, err)
		} else {
			d0 = append(d0, string(d))
		}
	}

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		for j := range d0 {
			pr, err := getParser(gotmLang, d0[j])
			if err != nil {
				b.Fatal(err)
				return
			}
			_, err = pr.Parse()
			if err != nil {
				b.Fatal(err)
				return
			}
		}
	}
}

func getParser(fn string, data string) (*Parser, error) {
	l, err := provider.GetLanguage(fn)
	if err != nil {
		return nil, err
	}
	return NewParser(l, []rune(data)), nil
}
