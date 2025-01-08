package main

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTopicStartsWith(t *testing.T) {
	tests := []struct {
		name     string
		topic    string
		prefix   string
		expected bool
	}{
		{"Valid Prefix", "sports_football", "sports", true},
		{"Invalid Prefix", "sports_basketball", "soccer", false},
		{"Exact Match", "sports", "sports", true},
		{"Empty Topic", "", "sports", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data := &templateData{
				Topic: tt.topic,
			}
			result := data.TopicStartsWith(tt.prefix)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestNewTemplateCache(t *testing.T) {
	err := os.MkdirAll("./ui/html/pages", os.ModePerm)
	assert.NoError(t, err, "Failed to create fake templates directory")
	err = os.MkdirAll("./ui/html/partials", os.ModePerm)
	assert.NoError(t, err, "Failed to create fake partials directory")
	defer os.RemoveAll("./ui")

	templateContent := `{{define "content"}}
		<h1>{{.Title}}</h1>
  {{end}}`

	mockBaseTemplate := `{{define "base"}}
<html>
	<head><title>Test Base Template</title></head>
	<body>
    {{template "content" .}}
		{{template "nav" .}}
	</body>
</html>
  {{end}}`

	mockNavTemplate := `{{define "nav"}}<nav>Navigation Bar</nav>{{end}}`

	templateFile := filepath.Join("./ui/html/pages", "test.tmpl")
	err = os.WriteFile(templateFile, []byte(templateContent), 0644)
	assert.NoError(t, err, "Failed to create fake template file")

	baseTemplateFile := filepath.Join("./ui/html", "base.tmpl")
	err = os.WriteFile(baseTemplateFile, []byte(mockBaseTemplate), 0644)
	assert.NoError(t, err, "Failed to create fake base template file")

	navTemplateFile := filepath.Join("./ui/html/partials", "nav.tmpl")
	err = os.WriteFile(navTemplateFile, []byte(mockNavTemplate), 0644)
	assert.NoError(t, err, "Failed to create fake nav template file")

	cache, err := newTemplateCache()
	assert.NoError(t, err)

	parsedTemplate, ok := cache["test.tmpl"]
	assert.True(t, ok, "Expected template to be in cache")

	var result bytes.Buffer

	err = parsedTemplate.ExecuteTemplate(&result, "base", map[string]string{"Title": "Test Page"})
	assert.NoError(t, err)

	output := result.String()
	assert.Contains(t, output, "Test Page")
	assert.Contains(t, output, "<h1>Test Page</h1>")
	assert.Contains(t, output, "Navigation Bar")
}
