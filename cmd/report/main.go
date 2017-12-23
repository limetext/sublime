// Copyright 2016 The lime Authors.
// Use of this source code is governed by a 2-clause
// BSD-style license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"log"
	"sort"

	"github.com/limetext/gopy"
	_ "github.com/limetext/sublime/api"
	"github.com/limetext/sublime/internal/util"
)

func main() {
	l := py.NewLock()
	defer l.Unlock()
	subl, err := py.Import("sublime")
	if err != nil {
		log.Fatalf("Error on importing sublime: %s", err)
	}

	sup := make(map[string][]string)
	exp := make(map[string][]string)
	if err := util.ExtractAPI(subl, sup, ""); err != nil {
		log.Fatalf("Error on extracting api: %s", err)
	} else if err := util.ReadAPI("data/api", exp); err != nil {
		log.Fatalf("Error reading data/api to api: %s", err)
	}

	keys := make([]string, 0)
	for key, _ := range exp {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	for _, key := range keys {
		exps := exp[key]
		sups, ok := sup[key]

		missings := make([]string, 0)
		if !ok {
			missings = append(missings, exps...)
			fmt.Printf("%s: 0%%\n", key)
			printMissing(missings)
			continue
		} else if len(exps) == 0 {
			fmt.Printf("%s: 100%%\n\n", key)
			continue
		}
		var count float64
		base := float64(len(exps))
		for _, v := range exps {
			if util.Exists(v, sups) {
				count++
			} else {
				missings = append(missings, v)
			}
		}
		fmt.Printf("%s: %d%%\n", key, int((count/base)*100))
		printMissing(missings)
	}
}

func printMissing(missings []string) {
	l := len(missings)
	if l == 0 {
		fmt.Print("\n")
		return
	}

	fmt.Print("\tMissings:")
	for i, missing := range missings {
		fmt.Printf(" %s", missing)
		if i != l-1 {
			fmt.Print(",")
		}
	}
	fmt.Print("\n\n")
}
