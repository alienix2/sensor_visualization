package main

import (
	"net/http"

	"github.com/alienix2/sensor_visualization/internal/validator"
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type userLoginForm struct {
	validator.Validator `form:"validator"`
	Email               string `form:"email"`
	Password            string `form:"password"`
}

func (app *application) userLogin(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
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
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, http.StatusUnprocessableEntity, "login.tmpl", data)
		return
	}

	username := app.accounts.GetUsername(id)
	mqtt_client, err := app.ConnectToBroker(username, form.Password)
	if err != nil {
		form.AddError("Global", "Unable to connect to MQTT broker, pls try again later")
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, http.StatusUnprocessableEntity, "login.tmpl", data)
		return
	}

	if client, ok := app.mqttClients.Load(id); ok {
		client.(mqtt.Client).Disconnect(250)
		app.mqttClients.Delete(id)
	}

	app.mqttClients.Store(id, mqtt_client)

	err = app.sessionManager.RenewToken(r.Context())
	if err != nil {
		app.errorLog.Println(err.Error())
		http.Error(w, "Internal server error", 500)
		return
	}
	app.sessionManager.Put(r.Context(), "authenticatedUserID", id)

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (app *application) userLogout(w http.ResponseWriter, r *http.Request) {
	err := app.sessionManager.RenewToken(r.Context())
	if err != nil {
		app.errorLog.Println(err.Error())
		http.Error(w, "Internal server error", 500)
		return
	}

	id := app.sessionManager.GetInt(r.Context(), "authenticatedUserID")
	if client, ok := app.mqttClients.Load(id); ok {
		client.(mqtt.Client).Disconnect(250)
		app.mqttClients.Delete(id)
	}
	app.sessionManager.Remove(r.Context(), "authenticatedUserID")
	http.Redirect(w, r, "/", http.StatusSeeOther)

	app.sessionManager.Remove(r.Context(), "authenticatedUserID")
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
