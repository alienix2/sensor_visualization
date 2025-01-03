package main

import (
	"fmt"
	"net/http"

	"github.com/alienix2/sensor_visualization/internal/validator"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	topics, err := app.topics.GetTopics()
	if err != nil {
		app.errorLog.Println(err.Error())
	}

	data := &templateData{
		Topics: topics,
	}
	app.render(w, http.StatusOK, "home.tmpl", data)
}

func (app *application) messagesByTopic(w http.ResponseWriter, r *http.Request) {
	topic := r.URL.Query().Get("topic")
	if topic == "" {
		app.errorLog.Println("Topic is required")
		http.Error(w, "Internal server error", 500)
		return
	}

	messages, err := app.topics.GetMessagesByTopic(topic)
	if err != nil {
		app.errorLog.Println(err.Error())
		http.Error(w, "Unable to fetch messages", http.StatusInternalServerError)
		return
	}

	data := &templateData{
		MessageData: messages,
		Topic:       topic,
	}

	app.render(w, http.StatusOK, "message_data.tmpl", data)
}

type userLoginForm struct {
	validator.Validator `form:"validator"`
	Email               string `form:"email"`
	Password            string `form:"password"`
}

func (app *application) userLogin(w http.ResponseWriter, r *http.Request) {
	data := &templateData{}
	data.Form = userLoginForm{}

	app.render(w, http.StatusOK, "login.tmpl", data)
}

func (app *application) userLoginPost(w http.ResponseWriter, r *http.Request) {
	var form userLoginForm
	form.Email = r.FormValue("email")
	form.Password = r.FormValue("password")

	id, err := app.accounts.Authenticate(form.Email, form.Password)
	if err != nil {
		form.AddError("EmailPass", "Email or password is incorrect")
		data := &templateData{}
		data.Form = form
		app.render(w, http.StatusUnprocessableEntity, "login.tmpl", data)
		return
	}

	fmt.Fprintf(w, "Authenticated with id: %d", id)
}

func (app *application) sensorsAdd(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("sensors add ToDo"))
}
