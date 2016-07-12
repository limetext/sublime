// Copyright 2013 The lime Authors.
// Use of this source code is governed by a 2-clause
// BSD-style license that can be found in the LICENSE file.

package theme

import (
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/limetext/loaders"
	"github.com/limetext/util"
)

func TestLoad(t *testing.T) {
	type Test struct {
		in  string
		out string
	}
	tests := []Test{
		{
			"../../../testdata/package/Monokai.tmTheme",
			"testdata/Monokai.tmTheme.res",
		},
	}
	for _, test := range tests {
		if d, err := ioutil.ReadFile(test.in); err != nil {
			t.Logf("Couldn't load file %s: %s", test.in, err)
		} else {
			var theme Theme
			if err := loaders.LoadPlist(d, &theme); err != nil {
				t.Error(err)
			} else {
				str := fmt.Sprintf("%s", theme)
				if d, err := ioutil.ReadFile(test.out); err != nil {
					if err := ioutil.WriteFile(test.out, []byte(str), 0644); err != nil {
						t.Error(err)
					}
				} else if diff := util.Diff(string(d), str); diff != "" {
					t.Error(diff)
				}

			}
		}
	}
}

func TestLoadFromPlist(t *testing.T) {
	f := "../../../testdata/package/Monokai.tmTheme"
	th, err := Load(f)
	if err != nil {
		t.Errorf("Tried to load %s, but got an error: %v", f, err)
	}

	n := "Monokai"
	if th.Name != n {
		t.Errorf("Tried to load %s, but got %s", f, th)
	}
}

func TestLoadFromNonPlist(t *testing.T) {
	f := "testdata/Monokai.tmTheme.res"
	_, err := Load(f)
	if err == nil {
		t.Errorf("Tried to load %s, expecting an error, but didn't get one", f)
	}
}

func TestLoadFromMissingFile(t *testing.T) {
	f := "testdata/MissingFile"
	_, err := Load(f)
	if err == nil {
		t.Errorf("Tried to load %s, expecting an error, but didn't get one", f)
	}
}

func TestGlobal(t *testing.T) {
	f := "../../../testdata/package/Monokai.tmTheme"
	th, err := Load(f)
	if err != nil {
		t.Fatalf("Tried to load %s, but got an error: %v", f, err)
	}
	gb := th.Settings()
	def := th.ScopeSettings[0].Settings
	if got, exp := gb.Background, def["background"]; got != exp {
		t.Errorf("Expected global settings background %s, but got %s", exp, got)
	}
	if got, exp := gb.Caret, def["caret"]; got != exp {
		t.Errorf("Expected global settings caret %s, but got %s", exp, got)
	}
	if got, exp := gb.LineHighlight, def["lineHighlight"]; got != exp {
		t.Errorf("Expected global settings lineHighlight %s, but got %s", exp, got)
	}
	if got, exp := gb.Selection, def["selection"]; got != exp {
		t.Errorf("Expected global settings selection %s, but got %s", exp, got)
	}
}
