// Copyright 2013 The lime Authors.
// Use of this source code is governed by a 2-clause
// BSD-style license that can be found in the LICENSE file.

package theme

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"sort"
	"strings"

	"github.com/limetext/backend/log"
	"github.com/limetext/backend/render"
	"github.com/limetext/loaders"
	"github.com/limetext/util"
)

type (
	// For loading tmTheme files
	Theme struct {
		Name     string
		Settings []ScopeSetting
		UUID     string
	}

	ScopeSetting struct {
		Name     string
		Scope    string
		Settings Settings
	}

	// TODO(q): personally I don't care about the font style attributes
	Settings map[string]render.Colour
)

func Load(filename string) (*Theme, error) {
	var scheme Theme
	if d, err := ioutil.ReadFile(filename); err != nil {
		return nil, fmt.Errorf("Unable to read theme definition: %s", err)
	} else if err := loaders.LoadPlist(d, &scheme); err != nil {
		return nil, fmt.Errorf("Unable to load theme definition: %s", err)
	}

	return &scheme, nil
}

func (s ScopeSetting) String() (ret string) {
	ret = fmt.Sprintf("%s - %s\n", s.Name, s.Scope)
	keys := make([]string, 0, len(s.Settings))
	for k := range s.Settings {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		ret += fmt.Sprintf("\t\t%s: %s\n", k, s.Settings[k])
	}
	return
}

func (t Theme) String() (ret string) {
	ret = fmt.Sprintf("%s - %s\n", t.Name, t.UUID)
	for i := range t.Settings {
		ret += fmt.Sprintf("\t%s", t.Settings[i])
	}
	return
}

func (s *Settings) UnmarshalJSON(data []byte) error {
	*s = make(Settings)
	tmp := make(map[string]json.RawMessage)
	if err := json.Unmarshal(data, &tmp); err != nil {
		return err
	}
	for k, v := range tmp {
		if strings.HasPrefix(k, "font") {
			continue
		}
		var c render.Colour
		if err := json.Unmarshal(v, &c); err != nil {
			return err
		}
		(*s)[k] = c
	}
	return nil
}

func (t *Theme) ClosestMatchingSetting(scope string) *ScopeSetting {
	pe := util.Prof.Enter("ClosestMatchingSetting")
	defer pe.Exit()
	na := scope
	for len(na) > 0 {
		sn := na
		i := strings.LastIndex(sn, " ")
		if i != -1 {
			sn = sn[i+1:]
		}

		for j := range t.Settings {
			if t.Settings[j].Scope == sn {
				return &t.Settings[j]
			}
		}
		if i2 := strings.LastIndex(na, "."); i2 == -1 {
			break
		} else if i > i2 {
			na = na[:i]
		} else {
			na = strings.TrimSpace(na[:i2])
		}
	}
	return &t.Settings[0]
}

func (t *Theme) Spice(vr *render.ViewRegions) (ret render.Flavour) {
	pe := util.Prof.Enter("Spice")
	defer pe.Exit()
	if len(t.Settings) == 0 {
		return
	}
	// If the scope hadn't wanted setting we use from global settings
	def := &t.Settings[0]

	s := t.ClosestMatchingSetting(vr.Scope)
	fg, ok := s.Settings["foreground"]
	if !ok {
		fg = def.Settings["foreground"]
	}
	ret.Foreground = render.Colour(fg)
	bname := "background"
	if vr.Flags&render.SELECTION != 0 {
		bname = "selection"
	}
	bg, ok := s.Settings[bname]
	if !ok {
		bg = def.Settings[bname]
	}
	ret.Background = render.Colour(bg)
	return
}

func (t *Theme) Global() (ret render.Global) {
	data, err := json.Marshal(t.Settings[0].Settings)
	if err != nil {
		log.Warn("Couldn't marshal global settings: %s", err)
		return
	}
	err = json.Unmarshal(data, &ret)
	if err != nil {
		log.Warn("Couldn't unmarshal to render.Global: %s", err)
	}
	return
}
