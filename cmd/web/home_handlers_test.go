package main

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	storage "github.com/alienix2/sensor_info/pkg/storage/central_database"
	"github.com/alienix2/sensor_visualization/internal/models"
	"github.com/alienix2/sensor_visualization/internal/models/mocks"
	"github.com/stretchr/testify/assert"
)

func TestHome(t *testing.T) {
	mockTopics := []*models.Topic{
		{ID: 1, Topic: "topic1"},
		{ID: 2, Topic: "topic2"},
	}

	tests := []struct {
		mockModel      *mocks.MockTopicModel
		name           string
		expectedStatus int
	}{
		{
			name: "ValidUser",
			mockModel: &mocks.MockTopicModel{
				Topics: mockTopics,
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "InvalidUser",
			mockModel: &mocks.MockTopicModel{
				Topics: mockTopics,
			},
			expectedStatus: http.StatusOK,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			app := newTestApplication()
			app.topics = test.mockModel

			ts := newTestServer(t, app.routes())
			defer ts.Close()

			code, _, _ := ts.get(t, "/")
			assert.Equal(t, test.expectedStatus, code)
		})
	}
}

func TestMessagesByTopic(t *testing.T) {
	mockMessages := []*storage.MessageData{
		{
			DeviceName: "Device 1",
			DeviceUnit: "Unit 1",
			DeviceData: 12.5,
			Topic:      "topic1",
			ID:         1,
			SentAt:     time.Now(),
		},
	}

	mockTopics := []*models.Topic{
		{ID: 1, Topic: "topic1"},
		{ID: 2, Topic: "topic2"},
	}

	mockModel := &mocks.MockTopicModel{
		Topics: mockTopics,
		Messages: map[string][]*storage.MessageData{
			"topic1": mockMessages,
		},
	}

	tests := []struct {
		mockModel      *mocks.MockTopicModel
		name           string
		topicQuery     string
		expectedBody   string
		userID         int
		expectedStatus int
	}{
		{
			name:           "ValidRequest",
			mockModel:      mockModel,
			userID:         1,
			topicQuery:     "topic1",
			expectedStatus: http.StatusOK,
			expectedBody:   "Device 1",
		},
		{
			name:           "InvalidTopic",
			mockModel:      mockModel,
			userID:         1,
			topicQuery:     "nonexistent_topic",
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   "Unable to fetch messages",
		},
		{
			name:           "InvalidUserID",
			mockModel:      mockModel,
			userID:         0,
			topicQuery:     "topic1",
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   "Unable to fetch messages",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			app := newTestApplication()
			app.topics = test.mockModel

			handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				app.messagesByTopic(w, r)
			})
			middleware := LoadAndSaveMock(app.sessionManager, "authenticatedUserID", test.userID)

			ts := newTestServer(t, middleware(handler))

			url := fmt.Sprintf("/messagesByTopic?topic=%s", test.topicQuery)
			resp, _, body := ts.get(t, url)

			assert.Equal(t, test.expectedStatus, resp)

			assert.Contains(t, body, test.expectedBody)
		})
	}
}
