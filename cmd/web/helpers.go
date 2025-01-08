package main

import (
	"fmt"
	"net/http"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/google/uuid"
)

func (app *application) render(w http.ResponseWriter, status int, page string, data *templateData) {
	ts, ok := app.templateCache[page]

	if !ok {
		app.errorLog.Println("Template not found " + page)
		http.Error(w, "Internal server error", 500)
		return
	}

	err := ts.ExecuteTemplate(w, "base", data)
	if err != nil {
		app.errorLog.Println(err.Error())
		http.Error(w, "Internal server error", 500)
	}
}

func (app *application) isAuthenticated(r *http.Request) bool {
	isAuthenticated := app.sessionManager.Exists(r.Context(), "authenticatedUserID")
	return isAuthenticated
}

func (app *application) ConnectToBroker(username, password string) (mqtt.Client, error) {
	opts := mqtt.NewClientOptions()
	opts.AddBroker(app.mqttBroker)
	opts.SetClientID(uuid.New().String())
	opts.SetUsername(username)
	opts.SetPassword(password)
	opts.SetTLSConfig(app.tlsConfig)

	mqttClient := mqtt.NewClient(opts)
	if token := mqttClient.Connect(); token.Wait() && token.Error() != nil {
		app.errorLog.Println(token.Error())
		return nil, token.Error()
	}
	return mqttClient, nil
}

func (app *application) sendCommandMessage(r *http.Request, topic, message string) error {
	id := app.sessionManager.GetInt(r.Context(), "authenticatedUserID")

	client, ok := app.mqttClients.Load(id)
	if !ok {
		app.errorLog.Println("Client not found")
		return fmt.Errorf("Client not found")
	}

	mqttClient := client.(mqtt.Client)
	token := mqttClient.Publish(topic, 1, true, message)
	token.Wait()

	return token.Error()
}
