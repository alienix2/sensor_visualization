package main

import (
	"log"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/alienix2/sensor_info/pkg/mqtt_utils"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/stretchr/testify/assert"
)

func TestControlMesageCreate_NotAuthenticated(t *testing.T) {
	app := newTestApplication()
	ts := newTestServer(t, app.routes())
	defer ts.Close()

	code, header, _ := ts.get(t, "/messages/sendcommand")
	assert.Equal(t, http.StatusSeeOther, code)
	assert.Equal(t, "/user/login", header.Get("Location"))
}

func TestControlMessageCreate(t *testing.T) {
	app := newTestApplication()

	handler := LoadAndSaveMock(app.sessionManager, "authenticatedUserID", 123)(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			app.ControlMessageCreatePost(w, r)
		}),
	)

	ts := newTestServer(t, handler)
	defer ts.Close()

	status, _, body := ts.get(t, "/messages/sendcommand")

	log.Println(body)
	assert.Equal(t, http.StatusOK, status)
	assert.Contains(t, body, "send_message.tmpl")
}

func TestControlMessageCreatePost(t *testing.T) {
	app := newTestApplication()

	brokerPort, err := mqtt_utils.GetAvailablePort()
	if err != nil {
		t.Fatalf("Error getting available port: %v", err)
	}

	mqtt_utils.StartMockMQTTServer(brokerPort)
	app.mqttBroker = "tcp://" + brokerPort

	opts := mqtt.NewClientOptions().AddBroker(app.mqttBroker)
	opts.SetClientID("testClient")

	mqttClient := mqtt.NewClient(opts)
	if token := mqttClient.Connect(); token.Wait() && token.Error() != nil {
		t.Fatalf("Error connecting to broker: %v", token.Error())
	}

	app.mqttClients.Store(123, mqttClient)

	handler := LoadAndSaveMock(app.sessionManager, "authenticatedUserID", 123)(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			app.ControlMessageCreatePost(w, r)
		}),
	)

	ts := newTestServer(t, handler)
	defer ts.Close()

	formData := url.Values{}
	formData.Add("topic", "test")
	formData.Add("message", `{"test": "test"}`)

	status, _, _ := ts.post(t, "/messages/sendcommand", strings.NewReader(formData.Encode()))

	assert.Equal(t, http.StatusSeeOther, status)
}

func TestControlMesageCreatePost_NotAuthenticated(t *testing.T) {
	app := newTestApplication()
	ts := newTestServer(t, app.routes())
	defer ts.Close()

	code, header, _ := ts.post(t, "/messages/sendcommand", nil)
	assert.Equal(t, http.StatusSeeOther, code)
	assert.Equal(t, "/user/login", header.Get("Location"))
}

func TestControlMessageCreatePost_WrongJSON(t *testing.T) {
	app := newTestApplication()
	app.mqttClients.Store(123, mqtt.NewClient(&mqtt.ClientOptions{}))

	handler := LoadAndSaveMock(app.sessionManager, "authenticatedUserID", 123)(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			app.ControlMessageCreatePost(w, r)
		}),
	)

	ts := newTestServer(t, handler)
	defer ts.Close()

	formData := url.Values{}
	formData.Add("topic", "test")
	formData.Add("message", `test`)

	status, _, body := ts.post(t, "/messages/sendcommand", strings.NewReader(formData.Encode()))

	assert.Equal(t, http.StatusOK, status)
	assert.Contains(t, body, "should be a valid JSON")
}

func TestControlMessageCreatePost_ClientDisconnected(t *testing.T) {
	app := newTestApplication()

	app.mqttClients.Store(123, mqtt.NewClient(&mqtt.ClientOptions{}))

	handler := LoadAndSaveMock(app.sessionManager, "authenticatedUserID", 123)(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			app.ControlMessageCreatePost(w, r)
		}),
	)

	ts := newTestServer(t, handler)
	defer ts.Close()

	form := url.Values{}
	form.Add("topic", "test")
	form.Add("message", `{"test": "test"}`)

	status, _, body := ts.post(t, "/messages/sendcommand", strings.NewReader(form.Encode()))

	assert.Equal(t, http.StatusOK, status)
	assert.Contains(t, body, "pls log-out and log-in again")
}
