// Copyright 2016 The lime Authors.
// Use of this source code is governed by a 2-clause
// BSD-style license that can be found in the LICENSE file.

package syntax

import (
	"io/ioutil"

	"gopkg.in/yaml.v1"

	"github.com/limetext/sublime/internal"
)

type (
	// for loading sublime-syntax files
	Syntax struct {
		Name           string
		FileExtensions []string `yaml:"file_extensions"`
		FirstLineMatch string   `yaml:"first_line_match"`
		Scope          string
		Hidden         bool
		Contexts       map[string][]Pattern
	}

	Pattern struct {
		Include              string
		MetaScope            string `yaml:"meta_scope"`
		MetaContentScope     string `yaml:"meta_content_scope"`
		MetaIncludePrototype string `yaml:"meta_include_prototype"`
		Match                internal.Regex
		Set                  string
		Scope                string
		Captures             internal.Captures
		Pop                  bool
		Push                 []Pattern
	}
)

func Load(filename string) (*Syntax, error) {
	var syn Syntax
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(data, &syn)
	if err != nil {
		return nil, err
	}

	return &syn, nil
}
