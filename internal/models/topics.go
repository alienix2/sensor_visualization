package models

import (
	storage "github.com/alienix2/sensor_info/pkg/storage/central_database"
	"gorm.io/gorm"
)

type TopicModelInterface interface {
	GetTopics(userID int) ([]*Topic, error)
	GetMessagesByTopic(userID int, topic string) ([]*storage.MessageData, error)
}

type Topic struct {
	Topic string
	ID    int
}

type TopicModel struct {
	DB *gorm.DB
}

func (t *TopicModel) GetTopics(userID int) ([]*Topic, error) {
	var topics []*Topic

	err := t.DB.Table("topics").
		Select("DISTINCT topics.*").
		Joins(`
			JOIN acl ON 
			(
				acl.topic NOT LIKE '#' AND topics.topic = acl.topic
			) 
			OR 
			(
				acl.topic LIKE '#' AND topics.topic LIKE REPLACE(acl.topic, '#', '%')
			)
		`).
		Where("acl.user_id = ? AND acl.rw IN (1, 3)", userID).
		Where("topics.topic NOT LIKE '%#'").
		Scan(&topics).Error
	if err != nil {
		return nil, err
	}

	return topics, nil
}

func (t *TopicModel) GetMessagesByTopic(userID int, topic string) ([]*storage.MessageData, error) {
	var messages []*storage.MessageData

	err := t.DB.Table("message_data").
		Select("message_data.*").
		Joins(`
			JOIN acl ON 
			(
				acl.topic NOT LIKE '#' AND message_data.topic = acl.topic
			) 
			OR 
			(
				acl.topic LIKE '#' AND message_data.topic LIKE REPLACE(acl.topic, '#', '%')
			)
		`).
		Where("acl.user_id = ? AND acl.rw IN (1, 3)", userID).
		Where("message_data.topic LIKE ?", topic).
		Order("message_data.created_at DESC").
		Scan(&messages).Error
	if err != nil {
		return nil, err
	}

	return messages, nil
}
