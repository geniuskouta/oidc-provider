package templates

import (
	"html/template"
	"net/http"
	"path/filepath"
)

// Engine is the struct that holds the template functions and templates map.
type Engine struct {
	templates *template.Template
}

// New creates a new Engine instance by parsing the templates from the provided directory.
func NewEngine(templatesDir string) (*Engine, error) {
	// Parse all templates in the specified directory.
	tmpl, err := template.ParseGlob(filepath.Join(templatesDir, "*.html"))
	if err != nil {
		return nil, err
	}

	return &Engine{templates: tmpl}, nil
}

// Render renders a template with the provided data and writes it to the ResponseWriter.
func (e *Engine) Render(w http.ResponseWriter, tmplName string, data interface{}) error {
	// Render the template with the data provided.
	return e.templates.ExecuteTemplate(w, tmplName, data)
}
