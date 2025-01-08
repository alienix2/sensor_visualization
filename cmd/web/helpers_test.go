package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	mock_broker "github.com/alienix2/sensor_info/pkg/mqtt_utils"
	"github.com/alienix2/sensor_visualization/cmd/web/mocks"
	"github.com/stretchr/testify/assert"
)

func TestConnectToBroker_Success(t *testing.T) {
	app := newTestApplication()

	port, err := mock_broker.GetAvailablePort()
	assert.NoError(t, err)

	mqttServer, err := mock_broker.StartMockMQTTServer(port)
	assert.NoError(t, err)
	defer mock_broker.StopMockMQTTServer(mqttServer)

	app.mqttBroker = "tcp://" + port

	client, err := app.ConnectToBroker("username", "password")
	assert.NoError(t, err)
	assert.NotNil(t, client)

	client.Disconnect(250)
}

func TestConnectToBroker_Failure(t *testing.T) {
	app := newTestApplication()

	app.mqttBroker = "tcp://localhost:12345"

	client, err := app.ConnectToBroker("username", "password")

	assert.Error(t, err)
	assert.Nil(t, client)
}

func TestSendCommandMessage(t *testing.T) {
	app := newTestApplication()

	port, err := mock_broker.GetAvailablePort()
	assert.NoError(t, err)

	mqttServer, err := mock_broker.StartMockMQTTServer(port)
	assert.NoError(t, err)
	defer mock_broker.StopMockMQTTServer(mqttServer)

	app.mqttBroker = "tcp://" + port

	client, err := app.ConnectToBroker("username", "password")
	assert.NoError(t, err)

	app.mqttClients.Store(123, client)
	handler := LoadAndSaveMock(app.sessionManager, "authenticatedUserID", 123)(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			err := app.sendCommandMessage(r, "test/topic", "test message")
			if err != nil {
				t.Fatalf("Error sending command message: %v", err)
			}
			w.WriteHeader(http.StatusOK)
		}),
	)

	ts := newTestServer(t, handler)
	defer ts.Close()

	status, _, _ := ts.post(t, "/messages/sendcommand", nil)

	assert.Equal(t, http.StatusOK, status)
}

func TestRender(t *testing.T) {
	app := newTestApplication()

	rec := httptest.NewRecorder()

	data := &templateData{
		Topic: "test",
	}
	app.render(rec, 200, "home.tmpl", data)

	assert.Contains(t, rec.Body.String(), "<p>test</p>")
	assert.Contains(t, rec.Body.String(), "<title>My App</title>")
	assert.Contains(t, rec.Body.String(), "<nav>")
	assert.Contains(t, rec.Body.String(), "<h1>home.tmpl</h1>")
}

func TestRenderTemplateNotFound(t *testing.T) {
	app := newTestApplication()

	app.templateCache = mocks.FakeTemplateCache()
	rec := httptest.NewRecorder()
	app.render(rec, http.StatusInternalServerError, "nonexistent.tmpl", nil)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}
