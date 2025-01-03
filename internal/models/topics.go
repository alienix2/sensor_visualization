package models

import (
	"strings"

	storage "github.com/alienix2/sensor_info/pkg/storage/central_database"
	"gorm.io/gorm"
)

type Topic struct {
	Topic string
	ID    int
}

type TopicModel struct {
	DB *gorm.DB
}

func (t *TopicModel) GetTopics() ([]*Topic, error) {
	var topics []*Topic
	err := t.DB.Find(&topics).Error
	if err != nil {
		return nil, err
	}

	return topics, nil
}

func (t *TopicModel) GetMessagesByTopic(topic string) ([]*storage.MessageData, error) {
	var messages []*storage.MessageData

	filter := strings.ReplaceAll(topic, "#", "%")

	err := t.DB.Where("topic LIKE ?", filter).Find(&messages).Error
	if err != nil {
		return nil, err
	}

	return messages, nil
}
