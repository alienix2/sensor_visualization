package mocks

import (
	"errors"
	"log"

	storage "github.com/alienix2/sensor_info/pkg/storage/central_database"
	models "github.com/alienix2/sensor_visualization/internal/models"
)

type MockTopicModel struct {
	Messages map[string][]*storage.MessageData
	Topics   []*models.Topic
}

func (m *MockTopicModel) GetTopics(userID int) ([]*models.Topic, error) {
	if userID == 0 {
		return nil, errors.New("invalid user ID")
	}
	return m.Topics, nil
}

func (m *MockTopicModel) GetMessagesByTopic(userID int, topic string) ([]*storage.MessageData, error) {
	if userID == 0 {
		return nil, errors.New("invalid user ID")
	}
	log.Println(m.Messages, topic)
	log.Println(m.Messages[topic])
	if messages, exists := m.Messages[topic]; exists {
		return messages, nil
	}
	return nil, errors.New("no messages found for the given topic")
}
