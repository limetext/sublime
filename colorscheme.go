// Copyright 2016 The lime Authors.
// Use of this source code is governed by a 2-clause
// BSD-style license that can be found in the LICENSE file.

package sublime

import (
	"path/filepath"

	"github.com/limetext/sublime/textmate/theme"
)

// wrapper arount Theme implements backend.ColorScheme
type colorScheme struct {
	*theme.Theme
}

func newColorScheme(path string) (*colorScheme, error) {
	if tm, err := theme.Load(path); err != nil {
		return nil, err
	} else {
		return &colorScheme{tm}, nil
	}
}

func (c *colorScheme) Name() string {
	return c.Theme.Name
}

func isColorScheme(path string) bool {
	if filepath.Ext(path) == ".tmTheme" {
		return true
	}
	return false
}
