package main

import (
	"html/template"
	"net/http"
	"path/filepath"
	"strings"

	storage "github.com/alienix2/sensor_info/pkg/storage/central_database"
	"github.com/alienix2/sensor_visualization/internal/models"
)

type templateData struct {
	Form            any
	Topic           string
	Topics          []*models.Topic
	MessageData     []*storage.MessageData
	IsAuthenticated bool
}

func newTemplateCache() (map[string]*template.Template, error) {
	cache := make(map[string]*template.Template)

	pages, err := filepath.Glob("./ui/html/pages/*.tmpl")
	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		name := filepath.Base(page)

		files := []string{
			"./ui/html/base.tmpl",
			"./ui/html/partials/nav.tmpl",
			page,
		}

		ts, err := template.ParseFiles(files...)
		if err != nil {
			return nil, err
		}

		cache[name] = ts
	}
	return cache, nil
}

func (app *application) newTemplateData(r *http.Request) *templateData {
	return &templateData{
		IsAuthenticated: app.isAuthenticated(r),
	}
}

func (d *templateData) TopicStartsWith(prefix string) bool {
	return strings.HasPrefix(d.Topic, prefix)
}
