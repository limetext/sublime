// Copyright 2016 The lime Authors.
// Use of this source code is governed by a 2-clause
// BSD-style license that can be found in the LICENSE file.

package syntax

import (
	"fmt"
	"testing"

	"github.com/gobs/pretty"
)

func Test(t *testing.T) {
	syn, err := Load("testdata/Go.sublime-syntax")
	if err != nil {
		t.Error(err)
	}
	fmt.Printf("%+v\n", pretty.PrettyFormat(syn))
}
