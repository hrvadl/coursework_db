package templates

import (
	"fmt"
	"html/template"
	"io/fs"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/hrvadl/coursework_db/pkg/config"
)

var htmlExt = ".html"

func parseTemplatesRecursive(path string) (*template.Template, error) {
	parsed := make(map[string]struct{})

	err := filepath.Walk(path, func(path string, info fs.FileInfo, err error) error {
		if info.IsDir() || !strings.HasSuffix(path, htmlExt) {
			return nil
		}

		if _, already := parsed[path]; already {
			return fmt.Errorf("you cannot have multiple templates with the same name: %v", path)
		}

		parsed[path] = struct{}{}
		return nil
	})

	if err != nil {
		return nil, err
	}

	var paths []string
	for path := range parsed {
		paths = append(paths, path)
	}

	return template.ParseFiles(paths...)
}

func NewResolver(cfg *config.Config) *Resolver {
	tmpl := template.Must(
		parseTemplatesRecursive(cfg.TemplatesDir),
	)

	return &Resolver{
		cfg:       cfg,
		templates: tmpl,
	}
}

type Resolver struct {
	cfg       *config.Config
	templates *template.Template
}

func (r *Resolver) Execute(w http.ResponseWriter, name string, data any) error {
	return r.templates.ExecuteTemplate(w, name, data)
}
