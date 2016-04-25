// Copyright 2016 The lime Authors.
// Use of this source code is governed by a 2-clause
// BSD-style license that can be found in the LICENSE file.

package internal

import (
	"encoding/json"
	"sort"
	"strconv"
)

type (
	Captures []Capture

	Capture struct {
		Key int
		Named
	}

	Named struct {
		Name string
	}
)

func (c *Captures) UnmarshalJSON(data []byte) error {
	tmp := make(map[string]Named)
	if err := json.Unmarshal(data, &tmp); err != nil {
		return err
	}
	for k, v := range tmp {
		i, _ := strconv.ParseInt(k, 10, 32)
		*c = append(*c, Capture{Key: int(i), Named: v})
	}
	sort.Sort(c)
	return nil
}

func (c *Captures) SetYAML(tag string, value interface{}) bool {
	tmp, ok := value.(map[interface{}]interface{})
	if !ok {
		return false
	}
	for k, v := range tmp {
		*c = append(*c, Capture{Key: k.(int), Named: Named{Name: v.(string)}})
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

func (c *Captures) Copy() *Captures {
	ret := make(Captures, len(*c))
	copy(ret, *c)
	return &ret
}
