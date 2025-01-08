package mocks

import (
	"html/template"
	"path/filepath"
)

func FakeTemplateCache() map[string]*template.Template {
	cache := make(map[string]*template.Template)

	base := `
		{{define "base"}}
		<!DOCTYPE html>
		<html lang="en">
		<head>
			<meta charset="UTF-8">
			<title>{{block "title" .}}My App{{end}}</title>
		</head>
		<body>
			<header>{{template "nav" .}}</header>
			<main>{{block "content" .}}{{end}}</main>
      {{if .MessageData}}
        {{range .MessageData}}
          <p>{{.}}</p>
        {{end}}
      {{end}}
          {{if .Topic}}
            <p>{{.Topic}}</p>
          {{end}}
          {{if .Form.Errors.EmailPass}}
            {{.Form.Errors.EmailPass}}
          {{end}}
          {{if .Form.Errors.Global}}
            {{.Form.Errors.Global}}
          {{end}}
          {{if .Form.Errors.Message}}
            {{.Form.Errors.Message}}
          {{end}}
		</body>
		</html>
		{{end}}
	`

	nav := `
		{{define "nav"}}
		<nav>
			<ul>
				<li><a href="/">Home</a></li>
				<li><a href="/user/login">Login</a></li>
				<li><a href="/messages/sendcommand">Send Message</a></li>
			</ul>
		</nav>
		{{end}}
	`

	pages := []string{"home.tmpl", "login.tmpl", "message_data.tmpl", "send_message.tmpl"}
	for _, page := range pages {
		content := `
			{{define "content"}}
			<h1>` + filepath.Base(page) + `</h1>
			<p>Welcome to the ` + filepath.Base(page) + ` page.</p>
			{{end}}
		`

		ts, err := template.New(page).Parse(base)
		if err != nil {
			panic(err)
		}
		_, err = ts.Parse(nav)
		if err != nil {
			panic(err)
		}
		_, err = ts.Parse(content)
		if err != nil {
			panic(err)
		}

		cache[page] = ts
	}

	return cache
}
