package template

import (
	"github.com/oam-dev/kubevela/apis/types"
	"github.com/oam-dev/kubevela/references/plugins"
)

// Manager defines a manager for template
type Manager interface {
	IsTrait(key string) bool
	LoadTemplate(key string) (tmpl string)
}

// Load will load all installed capabilities and create a manager
func Load() (Manager, error) {
	caps, err := plugins.LoadAllInstalledCapability()
	if err != nil {
		return nil, err
	}
	m := newManager()
	for _, cap := range caps {
		t := &Template{}
		t.Captype = cap.Type
		t.Raw = cap.CueTemplate
		m.Templates[cap.Name] = t
	}
	return m, nil
}

// Template defines a raw template struct
type Template struct {
	Captype types.CapType
	Raw     string
}

type manager struct {
	Templates map[string]*Template
}

func newManager() *manager {
	return &manager{
		Templates: make(map[string]*Template),
	}
}

func (m *manager) IsTrait(key string) bool {
	t, ok := m.Templates[key]
	if !ok {
		return false
	}
	return t.Captype == types.TypeTrait
}

func (m *manager) LoadTemplate(key string) string {
	t, ok := m.Templates[key]
	if !ok {
		return ""
	}
	return t.Raw
}
