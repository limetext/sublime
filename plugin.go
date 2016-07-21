// Copyright 2016 The lime Authors.
// Use of this source code is governed by a 2-clause
// BSD-style license that can be found in the LICENSE file.

package sublime

import (
	"path/filepath"

	"github.com/limetext/backend/log"
	"github.com/limetext/backend/packages"
	"github.com/limetext/gopy"
)

// Sublime plugin which is a single python file
type plugin struct {
	path string
	name string
}

func newPlugin(fn string) packages.Package {
	return &plugin{path: fn}
}

// TODO: implement unload
func (p *plugin) Load() {
	dir, file := filepath.Split(p.Path())
	p.name = file
	name := filepath.Base(dir) + "." + file[:len(file)-3]
	s, err := py.NewUnicode(name)
	if err != nil {
		log.Warn(err)
		return
	}
	defer s.Decref()

	log.Debug("Loading plugin %s", name)
	l := py.NewLock()
	defer l.Unlock()
	if r, err := module.Base().CallMethodObjArgs("reload_plugin", s); err != nil {
		log.Warn(err)
		return
	} else if r != nil {
		r.Decref()
	}
}

func (p *plugin) UnLoad() {}

func (p *plugin) Name() string {
	return p.name
}

func (p *plugin) Path() string {
	return p.path
}

func (p *plugin) FileChanged(name string) {
	p.Load()
}

func isPlugin(filename string) bool {
	return filepath.Ext(filename) == ".py"
}

var (
	pluginRecord = &packages.Record{isPlugin, newPlugin}

	module *py.Module
)

func pyAddPath(p string) {
	l := py.NewLock()
	defer l.Unlock()
	py.AddToPath(p)
}

func pyImport(name string) (*py.Module, error) {
	l := py.NewLock()
	defer l.Unlock()
	return py.Import(name)
}
