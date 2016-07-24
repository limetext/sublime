// Copyright 2013 The lime Authors.
// Use of this source code is governed by a 2-clause
// BSD-style license that can be found in the LICENSE file.

package api

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"

	"github.com/limetext/backend"
	_ "github.com/limetext/commands"
	"github.com/limetext/gopy"
	"github.com/limetext/util"
)

const (
	apifile = "testdata/api"
	supfile = "testdata/supported_api"
)

// Checking if we added necessary exported functions to sublime module
func TestSublimeApiMatch(t *testing.T) {
	l := py.NewLock()
	defer l.Unlock()
	subl, err := py.Import("sublime")
	if err != nil {
		t.Fatal(err)
	}

	sup := make(map[string][]string)
	if err := extractAPI(subl, sup, ""); err != nil {
		t.Fatalf("Error on extracting api: %s", err)
	}
	if exp, err := ioutil.ReadFile(supfile); err != nil {
		t.Fatalf("Error reading %s file: %s", supfile, err)
	} else if diff := util.Diff(string(exp), printMap(sup)); diff != "" {
		t.Errorf(diff)
	}

	exp := make(map[string][]string)
	if err := readAPI(apifile, exp); err != nil {
		t.Errorf("Error reading %s to api: %s", apifile, err)
	}
	skipKeys := []string{"sublime.TextCommandGlue", "sublime.ViewEventGlue", "sublime.ApplicationCommandGlue", "sublime.OnQueryContextGlue", "sublime.WindowCommandGlue"}
	skipVals := []string{"CLASS_CLOSING_PARENTHESIS", "CLASS_MIDDLE_WORD", "CLASS_OPENING_PARENTHESIS", "CLASS_WORD_END_WITH_PUNCTUATION", "CLASS_WORD_START_WITH_PUNCTUATION", "register", "unregister", "console"}
	for k, sups := range sup {
		if exists(k, skipKeys) {
			continue
		}
		exps, ok := exp[k]
		if !ok {
			t.Errorf("We have exported '%s', but its not exported in sublime text", k)
			continue
		}
		for _, v := range sups {
			if !exists(v, skipVals) && !exists(v, exps) {
				t.Errorf("We have exported %s:'%s', but its not exported in sublime text", k, v)
			}
		}
	}

	logSupPercent(sup, exp)
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

// Extracts the functionality we currently support to map
func extractAPI(v py.Object, m map[string][]string, key string) error {
	b := v.Base()
	dir, err := b.Dir()
	if err != nil {
		return err
	}
	defer dir.Decref()
	l, ok := dir.(*py.List)
	if !ok {
		return fmt.Errorf("Unexpected type: %v", dir.Type())
	}
	sl := l.Slice()
	for _, v2 := range sl {
		if str := fmt.Sprint(v2); strings.HasPrefix(str, "__") {
			continue
		}
		if key != "" {
			m[key] = append(m[key], fmt.Sprint(v2))
			continue
		}
		item, err := b.GetAttr(v2)
		if err != nil {
			return err
		}
		ty := item.Type()
		k := fmt.Sprint(v2)
		if k == "RegionSet" {
			k = "Selection"
		}
		if ty == py.TypeType {
			k = "sublime." + k
			m[k] = make([]string, 0)
			if err := extractAPI(item, m, k); err != nil {
				return err
			}
		} else {
			m["sublime"] = append(m["sublime"], fmt.Sprint(v2))
		}
		item.Decref()
	}
	return nil
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

func logSupPercent(sup, exp map[string][]string) {
	keys := make([]string, 0)
	for key, _ := range exp {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	for _, key := range keys {
		exps := exp[key]
		sups, ok := sup[key]
		if !ok {
			fmt.Printf("%s: 0%%\n", key)
			continue
		} else if len(exps) == 0 {
			fmt.Printf("%s: 100%%\n", key)
			continue
		}
		var count float64
		base := float64(len(exps))
		for _, v := range exps {
			if exists(v, sups) {
				count++
			}
		}
		fmt.Printf("%s: %d%%\n", key, int((count/base)*100))
	}
}

// Reads the api from file to map
func readAPI(fn string, m map[string][]string) error {
	f, err := os.Open(fn)
	if err != nil {
		return err
	}
	defer f.Close()

	scnr := bufio.NewScanner(f)
	var key string
	for scnr.Scan() {
		if s := scnr.Text(); strings.Contains(s, "//") {
			continue
		} else if strings.Contains(s, "\t") {
			m[key] = append(m[key], strings.Replace(s, "\t", "", -1))
			continue
		} else {
			m[s] = make([]string, 0)
			key = s
		}
	}
	return scnr.Err()
}

func exists(v string, ss []string) (exist bool) {
	for _, v2 := range ss {
		if v2 == v {
			return true
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
