package main

import (
	"encoding/json"
	"net/http"
	"net/url"

	"github.com/alienix2/sensor_visualization/internal/validator"
)

type commandForm struct {
	validator.Validator `form:"validator"`
	Topic               string `form:"topic"`
	Message             string `form:"message"`
}

func (app *application) ControlMessageCreate(w http.ResponseWriter, r *http.Request) {
	topic := r.URL.Query().Get("topic")

	if !app.isAuthenticated(r) {
		http.Redirect(w, r, "/user/login", http.StatusSeeOther)
		return
	}

	data := app.newTemplateData(r)
	data.Topic = topic
	data.Form = commandForm{}

	app.render(w, http.StatusOK, "send_message.tmpl", data)
}

func (app *application) ControlMessageCreatePost(w http.ResponseWriter, r *http.Request) {
	var form commandForm
	if !app.isAuthenticated(r) {
		http.Redirect(w, r, "/user/login", http.StatusSeeOther)
		return
	}

	form.Topic = r.FormValue("topic")
	form.Message = r.FormValue("message")

	form.CheckField(json.Valid([]byte(form.Message)), "Message", "This should be a valid JSON")

	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		data.Topic = form.Topic
		app.render(w, http.StatusUnprocessableEntity, "send_message.tmpl", data)
		return
	}

	err := app.sendCommandMessage(r, "command/"+form.Topic, form.Message)
	if err != nil {
		form.AddError("Global", "Unable to send message to broker, pls log-out and log-in again and check broker's connection")
		data := app.newTemplateData(r)
		data.Form = form
		data.Topic = form.Topic
		app.render(w, http.StatusUnprocessableEntity, "send_message.tmpl", data)
		return
	}

	http.Redirect(w, r, "/messages/view?topic="+url.QueryEscape(form.Topic), http.StatusSeeOther)
}
