package sublime

import (
	"io/ioutil"
	"os"
	"path"
	"path/filepath"

	"github.com/limetext/lime-backend/lib"
	"github.com/limetext/lime-backend/lib/keys"
	"github.com/limetext/lime-backend/lib/log"
	"github.com/limetext/lime-backend/lib/packages"
	"github.com/limetext/text"
)

type pkg struct {
	dir string
	text.HasSettings
	keys.HasKeyBindings
	platformSet *text.HasSettings
	defaultSet  *text.HasSettings
	defaultKB   *keys.HasKeyBindings
	plugins     map[string]*plugin
	// TODO: themes, snippets, etc more info on iss#71
}

func newPKG(dir string) packages.Package {
	p := &pkg{
		dir:         dir,
		platformSet: new(text.HasSettings),
		defaultSet:  new(text.HasSettings),
		defaultKB:   new(keys.HasKeyBindings),
		plugins:     make(map[string]*plugin),
	}

	return p
}

func (p *pkg) Load() {
	log.Debug("Loading package %s", p.Name())
	p.loadKeyBindings()
	p.loadSettings()
	p.loadPlugins()
}

func (p *pkg) Name() string {
	return p.dir
}

func (p *pkg) FileCreated(name string) {
	p.loadPlugin(name)
}

func (p *pkg) loadPlugins() {
	log.Fine("Loading %s plugins", p.Name())
	fis, err := ioutil.ReadDir(p.Name())
	if err != nil {
		log.Warn("Error on reading directory %s, %s", p.Name(), err)
		return
	}
	for _, fi := range fis {
		p.loadPlugin(fi.Name())
	}
}

func (p *pkg) loadPlugin(fn string) {
	if !isPlugin(fn) {
		return
	}
	_, exist := p.plugins[fn]
	if exist {
		return
	}

	pl := newPlugin(fn)
	pl.Load()
	packages.Watch(pl)

	p.plugins[fn] = pl.(*plugin)
}

func (p *pkg) loadKeyBindings() {
	log.Fine("Loading %s keybindings", p.Name())
	ed := backend.GetEditor()
	tmp := ed.KeyBindings().Parent()
	dir := filepath.Dir(p.Name())

	ed.KeyBindings().SetParent(p)
	p.KeyBindings().SetParent(p.defaultKB)
	p.defaultKB.KeyBindings().SetParent(tmp)

	pt := path.Join(dir, "Default.sublime-keymap")
	packages.NewKeymapL(pt, p.defaultKB.KeyBindings())

	pt = path.Join(dir, "Default ("+ed.Plat()+").sublime-keymap")
	packages.NewKeymapL(pt, p.KeyBindings())
}

func (p *pkg) loadSettings() {
	log.Fine("Loading %s settings", p.Name())
	ed := backend.GetEditor()
	tmp := ed.Settings().Parent()
	dir := filepath.Dir(p.Name())

	ed.Settings().SetParent(p)
	p.Settings().SetParent(p.platformSet)
	p.platformSet.Settings().SetParent(p.defaultSet)
	p.defaultSet.Settings().SetParent(tmp)

	pt := path.Join(dir, "Preferences.sublime-settings")
	packages.NewSettingL(pt, p.defaultSet.Settings())

	pt = path.Join(dir, "Preferences ("+ed.Plat()+").sublime-settings")
	packages.NewSettingL(pt, p.platformSet.Settings())

	pt = path.Join(ed.PackagesPath("user"), "Preferences.sublime-settings")
	packages.NewSettingL(pt, p.Settings())
}

func isPKG(dir string) bool {
	fm, err := os.Stat(dir)
	if err != nil {
		return false
	}
	return fm.IsDir()
}

func init() {
	packages.Register(packages.Record{isPKG, newPKG})
}