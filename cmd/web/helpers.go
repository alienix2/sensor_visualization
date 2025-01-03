package main

import (
	"net/http"
)

func (app *application) render(w http.ResponseWriter, status int, page string, data *templateData) {
	ts, ok := app.templateCache[page]

	if !ok {
		app.errorLog.Println("Template not found " + page)
		http.Error(w, "Internal server error", 500)
		return
	}

	w.WriteHeader(status)

	err := ts.ExecuteTemplate(w, "base", data)
	if err != nil {
		app.errorLog.Println(err.Error())
		http.Error(w, "Internal server error", 500)
	}
}
