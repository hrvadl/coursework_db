package templates

import (
	"path/filepath"

	"github.com/hrvadl/coursework_db/pkg/config"
)

func NewResolver(cfg *config.Config) *Resolver {
	return &Resolver{
		cfg: cfg,
	}
}

type Resolver struct {
	cfg *config.Config
}

func (r *Resolver) ResolveHTML(name string) string {
	return filepath.Join(r.cfg.TemplatesDir, name)
}
