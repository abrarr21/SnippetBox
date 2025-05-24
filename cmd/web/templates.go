package main

import (
	"html/template"
	"net/url"
	"path/filepath"
	"time"

	"github.com/abrarr21/snippet/pkg/models"
)

type templateData struct {
	CurrentYear int
	FormData    url.Values        // holds any previously submitted data
	FormErrors  map[string]string // holds any validation errors
	Snippet     *models.Snippet
	Snippets    []*models.Snippet
}

func humanDate(t time.Time) string {
	return t.Format("02 Jan 2006 at 15:04")
}

var functions = template.FuncMap{
	"humanDate": humanDate,
}

func newTemplateCache(dir string) (map[string]*template.Template, error) {
	// Initialize a new map to act as a cache.
	cache := map[string]*template.Template{}

	//Use the filepath.Glob function to get a slice of all the filepaths with the extension '.page.tmpl'. This essentially gives us a slice of all the 'page' templates for the application.
	pages, err := filepath.Glob(filepath.Join(dir, "*.page.html"))
	if err != nil {
		return nil, err
	}

	// Loop through the pages one by one.
	for _, page := range pages {
		//Extract the file name (like 'home.page.tmpl') from the full file path and assign it to the name variable.
		name := filepath.Base(page)

		// use the ParseGlob method to add any 'layout' template to the template set (in our case, it's just the 'base' layout at the moment.)
		ts, err := template.New(name).Funcs(functions).ParseFiles(page)
		if err != nil {
			return nil, err
		}

		// use the ParseGlob method to add any 'layout' template to the template set (in our case, it's just the 'footer' layout at the moment.)
		ts, err = ts.ParseGlob(filepath.Join(dir, "*.layout.html"))
		if err != nil {
			return nil, err
		}

		// use the ParseGlob method to add any 'partial' template to the template set (in our case, it's just the 'footer' layout at the moment.)
		ts, err = ts.ParseGlob(filepath.Join(dir, "*.partial.html"))
		if err != nil {
			return nil, err
		}

		// Add the template set to the cache, using the name of the page (like 'home.page.tmpl' ) as the key
		cache[name] = ts

	}
	// Return the map
	return cache, nil
}
