package main

import (
	"html/template"
	"path/filepath"

	storage "github.com/alienix2/sensor_info/pkg/storage/central_database"
	"github.com/alienix2/sensor_visualization/internal/models"
)

type templateData struct {
	Form        any
	Topic       string
	Topics      []*models.Topic
	MessageData []*storage.MessageData
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
