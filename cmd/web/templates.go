package main

import (
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"

	"github.com/jcorry/snippetbox/pkg/models"
)

type templateData struct {
	Snippet  *models.Snippet
	Snippets []*models.Snippet
}

// Renders template data in template from cache
func (app *application) render(w http.ResponseWriter, r *http.Request, name string, td *templateData) {
	// Get the template from cache
	ts, ok := app.templateCache[name]
	if !ok {
		app.serverError(w, fmt.Errorf("The template %s does not exist.", name))
		return
	}

	// execute the template
	err := ts.Execute(w, td)
	if err != nil {
		app.serverError(w, err)
	}
}

// Builds a cache of template files
func newTemplateCache(dir string) (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}

	pages, err := filepath.Glob(filepath.Join(dir, "*.page.tmpl"))
	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		name := filepath.Base(page)

		ts, err := template.ParseFiles(page)
		if err != nil {
			return nil, err
		}

		ts, err = ts.ParseGlob(filepath.Join(dir, "*.layout.tmpl"))
		if err != nil {
			return nil, err
		}

		ts, err = ts.ParseGlob(filepath.Join(dir, "*.partial.tmpl"))
		if err != nil {
			return nil, err
		}
		cache[name] = ts
	}

	fmt.Printf("%v", cache)
	return cache, nil
}
