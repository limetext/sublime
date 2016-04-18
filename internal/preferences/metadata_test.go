// Copyright 2016 The lime Authors.
// Use of this source code is governed by a 2-clause
// BSD-style license that can be found in the LICENSE file.

package preferences

import (
	"io/ioutil"
	"testing"

	"github.com/limetext/util"
)

func TestLoadMetadata(t *testing.T) {
	var (
		in  = "testdata/Comments.tmPreferences"
		exp = "testdata/Comments.tmPreferences.res"
	)

	md, err := LoadMetadata(in)
	if err != nil {
		t.Fatalf("Error on loading %s: %s", in, err)
	}
	data, err := ioutil.ReadFile(exp)
	if err != nil {
		t.Fatalf("Error reading expected file %s: %s", exp, err)
	}
	if diff := util.Diff(string(data), md.String()); diff != "" {
		t.Error(diff)
	}
}
