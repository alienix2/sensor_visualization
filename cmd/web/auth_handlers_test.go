package main

import (
	"net/http"
	"net/url"
	"strings"
	"sync"
	"testing"

	"github.com/alienix2/sensor_info/pkg/mqtt_utils"
	"github.com/alienix2/sensor_visualization/internal/models/mocks"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/stretchr/testify/assert"
)

func TestUserLogin(t *testing.T) {
	app := newTestApplication()
	ts := newTestServer(t, app.routes())
	defer ts.Close()

	status, _, body := ts.get(t, "/user/login")

	assert.Equal(t, http.StatusOK, status)
	assert.Contains(t, body, "Login")
}

func TestUserLoginPost_Success(t *testing.T) {
	mockEmail := "user@example.com"
	mockPassword := "correctPassword"
	mockUsername := "user123"

	app := newTestApplication()
	app.mqttClients = sync.Map{}

	brokerPort, err := mqtt_utils.GetAvailablePort()
	if err != nil {
		t.Fatalf("Error getting available port: %v", err)
	}
	mqtt_utils.StartMockMQTTServer(brokerPort)
	app.mqttBroker = "tcp://" + brokerPort

	mockAccountModel := &mocks.MockAccountModel{}
	mockAccountModel.Insert(mockUsername, mockEmail, mockPassword)
	app.accounts = mockAccountModel

	ts := newTestServer(t, app.routes())
	defer ts.Close()

	form := url.Values{
		"email":    {mockEmail},
		"password": {mockPassword},
	}
	status, headers, _ := ts.post(t, "/user/login", strings.NewReader(form.Encode()))

	assert.Equal(t, http.StatusSeeOther, status)
	assert.Equal(t, "/", headers.Get("Location"))
}

func TestUserLoginPost_BrokerUnreachable(t *testing.T) {
	mockEmail := "user@example.com"
	mockPassword := "correctPassword"
	mockUsername := "user123"

	app := newTestApplication()
	app.mqttClients = sync.Map{}

	mockAccountModel := &mocks.MockAccountModel{}
	mockAccountModel.Insert(mockUsername, mockEmail, mockPassword)
	app.accounts = mockAccountModel

	ts := newTestServer(t, app.routes())
	defer ts.Close()

	form := url.Values{
		"email":    {mockEmail},
		"password": {mockPassword},
	}
	status, _, body := ts.post(t, "/user/login", strings.NewReader(form.Encode()))

	assert.Equal(t, http.StatusOK, status)
	assert.Contains(t, body, "Unable to connect to MQTT broker")
}

func TestUserLoginPost_EmptyForm(t *testing.T) {
	mockEmail := "user@example.com"
	mockPassword := "correctPassword"
	mockUsername := "user123"

	app := newTestApplication()
	app.mqttClients = sync.Map{}

	mockAccountModel := &mocks.MockAccountModel{}
	mockAccountModel.Insert(mockUsername, mockEmail, mockPassword)
	app.accounts = mockAccountModel

	ts := newTestServer(t, app.routes())
	defer ts.Close()

	tests := []url.Values{
		{"email": {""}, "password": {mockPassword}},
		{"email": {mockEmail}, "password": {""}},
		{"email": {""}, "password": {""}},
	}

	for _, form := range tests {
		status, _, body := ts.post(t, "/user/login", strings.NewReader(form.Encode()))
		assert.Equal(t, http.StatusOK, status)
		assert.Contains(t, body, "Email or password is incorrect")
	}
}

func TestUserLogout(t *testing.T) {
	app := newTestApplication()
	app.mqttClients.Store(123, mqtt.NewClient(&mqtt.ClientOptions{}))

	handler := LoadAndSaveMock(app.sessionManager, "authenticatedUserID", 123)(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			app.userLogout(w, r)
		}),
	)
	ts := newTestServer(t, handler)
	defer ts.Close()

	status, headers, _ := ts.get(t, "/user/logout")

	assert.Equal(t, http.StatusSeeOther, status)
	assert.Equal(t, "/", headers.Get("Location"))

	client, ok := app.mqttClients.Load(123)
	assert.False(t, ok, "MQTT client should be removed")
	assert.Nil(t, client, "MQTT client should be nil after removal")
}
