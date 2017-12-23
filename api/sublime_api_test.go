// Copyright 2013 The lime Authors.
// Use of this source code is governed by a 2-clause
// BSD-style license that can be found in the LICENSE file.

package api

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"testing"

	"github.com/limetext/backend"
	_ "github.com/limetext/commands"
	"github.com/limetext/gopy"
	. "github.com/limetext/sublime/internal/util"
	"github.com/limetext/util"
)

const (
	apifile = "../data/api"
	supfile = "../data/supported_api"
)

// Check if we are exporting the expected api(report/supported_api)
func TestSublimeApiMatchBefore(t *testing.T) {
	l := py.NewLock()
	defer l.Unlock()
	subl, err := py.Import("sublime")
	if err != nil {
		t.Fatalf("Error on importing sublime: %s", err)
	}

	sup := make(map[string][]string)
	if err := ExtractAPI(subl, sup, ""); err != nil {
		t.Fatalf("Error on extracting api: %s", err)
	}
	if exp, err := ioutil.ReadFile(supfile); err != nil {
		t.Fatalf("Error reading %s file: %s", supfile, err)
	} else if diff := util.Diff(string(exp), printMap(sup)); diff != "" {
		t.Errorf(diff)
	}
}

// Check if we are exporting extra functionality
// All of exported api should exist in report/api
func TestExportedApi(t *testing.T) {
	skipKeys := []string{"sublime.TextCommandGlue", "sublime.ViewEventGlue", "sublime.ApplicationCommandGlue", "sublime.OnQueryContextGlue", "sublime.WindowCommandGlue"}
	skipVals := []string{"CLASS_CLOSING_PARENTHESIS", "CLASS_MIDDLE_WORD", "CLASS_OPENING_PARENTHESIS", "CLASS_WORD_END_WITH_PUNCTUATION", "CLASS_WORD_START_WITH_PUNCTUATION", "register", "unregister", "console"}

	l := py.NewLock()
	defer l.Unlock()
	subl, err := py.Import("sublime")
	if err != nil {
		t.Fatalf("Error on importing sublime: %s", err)
	}

	sup := make(map[string][]string)
	exp := make(map[string][]string)
	if err := ExtractAPI(subl, sup, ""); err != nil {
		t.Fatalf("Error on extracting api: %s", err)
	} else if err := ReadAPI(apifile, exp); err != nil {
		t.Fatalf("Error reading %s to api: %s", apifile, err)
	}

	for k, sups := range sup {
		if Exists(k, skipKeys) {
			continue
		}
		exps, ok := exp[k]
		if !ok {
			t.Errorf("We have exported '%s', but its not exported in sublime text", k)
			continue
		}
		for _, v := range sups {
			if !Exists(v, skipVals) && !Exists(v, exps) {
				t.Errorf("We have exported %s:'%s', but its not exported in sublime text", k, v)
			}
		}
	}
}

// basicly running "testdata/*.py" files
func TestSublimeApi(t *testing.T) {
	l := py.NewLock()
	defer l.Unlock()

	dir, err := os.Open("testdata")
	if err != nil {
		t.Fatal(err)
	}
	files, err := dir.Readdirnames(0)
	if err != nil {
		t.Fatal(err)
	}

	for _, fn := range files {
		if filepath.Ext(fn) != ".py" {
			continue
		}

		t.Logf("Running %s", fn)
		if _, err := py.Import(fn[:len(fn)-3]); err != nil {
			t.Error(err)
		} else {
			t.Logf("Ran %s", fn)
		}
	}
}

func printMap(m map[string][]string) (s string) {
	keys := make([]string, 0)
	for k, v := range m {
		keys = append(keys, k)
		sort.Strings(v)
	}
	sort.Strings(keys)
	for _, key := range keys {
		s += fmt.Sprintf("%s\n", key)
		ss := m[key]
		for _, v := range ss {
			s += fmt.Sprintf("\t%s\n", v)
		}
	}
	return
}

func init() {
	l := py.NewLock()
	defer l.Unlock()
	py.AddToPath("testdata")

	ed := backend.GetEditor()
	ed.Init()
	ed.NewWindow()
}
