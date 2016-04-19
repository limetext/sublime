// Copyright 2016 The lime Authors.
// Use of this source code is governed by a 2-clause
// BSD-style license that can be found in the LICENSE file.

package syntax

import (
	"io/ioutil"
	"sort"

	"github.com/limetext/sublime/internal/textmate"
	"gopkg.in/yaml.v1"
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
		Match                textmate.Regex
		Scope                string
		Captures             Captures
		Pop                  bool
		Push                 []Pattern
		Set                  []Pattern
	}

	Captures []Capture

	Capture struct {
		Key  int
		Name string
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

func (c *Captures) SetYAML(tag string, value interface{}) bool {
	tmp, ok := value.(map[interface{}]interface{})
	if !ok {
		return false
	}
	for k, v := range tmp {
		*c = append(*c, Capture{Key: k.(int), Name: v.(string)})
	}
	sort.Sort(c)
	return true
}

func (c *Captures) Len() int {
	return len(*c)
}

func (c *Captures) Less(i, j int) bool {
	return (*c)[i].Key < (*c)[j].Key
}

func (c *Captures) Swap(i, j int) {
	(*c)[i], (*c)[j] = (*c)[j], (*c)[i]
}

func (c *Captures) copy() *Captures {
	ret := make(Captures, len(*c))
	copy(ret, *c)
	return &ret
}
