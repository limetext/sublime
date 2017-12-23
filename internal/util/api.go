// Copyright 2016 The lime Authors.
// Use of this source code is governed by a 2-clause
// BSD-style license that can be found in the LICENSE file.

package util

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/limetext/gopy"
)

// Extracts python object api
func ExtractAPI(v py.Object, m map[string][]string, key string) error {
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
			if err := ExtractAPI(item, m, k); err != nil {
				return err
			}
		} else {
			m["sublime"] = append(m["sublime"], fmt.Sprint(v2))
		}
		item.Decref()
	}
	return nil
}

// Reads api from file to map
func ReadAPI(fn string, m map[string][]string) error {
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

func Exists(v string, ss []string) (exist bool) {
	for _, v2 := range ss {
		if v2 == v {
			return true
		}
	}
	return
}
