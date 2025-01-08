package models

import (
	"context"
	"testing"
	"time"

	storage "github.com/alienix2/sensor_info/pkg/storage/central_database"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type ACL struct {
	Topic  string `gorm:"primaryKey"`
	UserID int    `gorm:"primaryKey"`
	Rw     int
}

func (ACL) TableName() string {
	return "acl"
}

func setupTopicDatabase(t *testing.T) *gorm.DB {
	ctx := context.Background()

	dsn, err := storage.SetupContainer(ctx)
	if err != nil {
		t.Fatalf("Failed to set up container: %v", err)
	}

	err = storage.InitMySQLCentralDatabase(dsn)
	if err != nil {
		t.Fatalf("Failed to initialize the database: %v", err)
	}

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to connect to the database: %v", err)
	}

	if err := db.AutoMigrate(&Topic{}, &ACL{}, &storage.MessageData{}); err != nil {
		t.Fatalf("Failed to run migrations: %v", err)
	}

	return db
}

func TestTopicModel(t *testing.T) {
	db := setupTopicDatabase(t)
	topicModel := TopicModel{DB: db}

	t.Run("Test GetTopics", func(t *testing.T) {
		err := db.Create(&Topic{Topic: "test_topic"}).Error
		assert.NoError(t, err)

		err = db.Create(&ACL{
			UserID: 1, Topic: "test_topic", Rw: 1,
		}).Error
		assert.NoError(t, err)

		topics, err := topicModel.GetTopics(1)
		assert.NoError(t, err)
		assert.NotEmpty(t, topics)
		assert.Equal(t, "test_topic", topics[0].Topic)
	})

	t.Run("Test GetMessagesByTopic", func(t *testing.T) {
		err := db.Create(&storage.MessageData{
			Topic:      "test_topic",
			DeviceName: "Device1",
			DeviceUnit: "Celsius",
			DeviceID:   "device-123",
			Notes:      "Test message 1",
			DeviceData: 25.5,
			SentAt:     time.Now(),
		}).Error
		assert.NoError(t, err)

		messages, err := topicModel.GetMessagesByTopic(1, "test_topic")
		assert.NoError(t, err)
		assert.NotEmpty(t, messages)
		assert.Equal(t, "Device1", messages[0].DeviceName)
		assert.Equal(t, "Celsius", messages[0].DeviceUnit)
		assert.Equal(t, 25.5, messages[0].DeviceData)
	})
}
