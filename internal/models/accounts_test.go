package models

import (
	"context"
	"testing"

	central_storage "github.com/alienix2/sensor_info/pkg/storage/central_database"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func TestAccountModel(t *testing.T) {
	ctx := context.Background()

	dsn, err := central_storage.SetupContainer(ctx)
	if err != nil {
		t.Fatalf("Failed to set up container: %v", err)
	}
	defer central_storage.TerminateMariaDBContainer(ctx)

	err = central_storage.InitMySQLCentralDatabase(dsn)
	if err != nil {
		t.Fatalf("Failed to initialize the database: %v", err)
	}

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to connect to the database: %v", err)
	}
	accountModel := AccountModel{DB: db}

	if err := db.AutoMigrate(Account{}); err != nil {
		t.Fatalf("Failed to run migrations: %v", err)
	}

	t.Run("Test Insert and Authenticate", func(t *testing.T) {
		err := accountModel.Insert("testuser", "testuser@example.com", "password123")
		assert.NoError(t, err)

		accountID, err := accountModel.Authenticate("testuser@example.com", "password123")
		assert.NoError(t, err)
		assert.NotEqual(t, accountID, 0)
	})

	t.Run("Test GetUsername", func(t *testing.T) {
		err := accountModel.Insert("testuser2", "testuser2@example.com", "password123")
		assert.NoError(t, err)

		accountID, err := accountModel.Authenticate("testuser2@example.com", "password123")
		assert.NoError(t, err)

		username := accountModel.GetUsername(accountID)
		assert.Equal(t, "testuser2", username)
	})
}
