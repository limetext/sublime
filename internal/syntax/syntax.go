// Copyright 2016 The lime Authors.
// Use of this source code is governed by a 2-clause
// BSD-style license that can be found in the LICENSE file.

package syntax

import (
	"io/ioutil"

	"gopkg.in/yaml.v1"
)

// For loading sublime-syntax files
// https://www.sublimetext.com/docs/3/syntax.html
type Syntax struct {
	Name           string
	FileExtensions []string `yaml:"file_extensions"`
	FirstLineMatch string   `yaml:"first_line_match"`
	Scope          string
	Variables      map[string]string
	Hidden         bool
	Contexts       map[string]Context
}

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
