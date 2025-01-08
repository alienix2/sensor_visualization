package main

import (
	"context"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/alexedwards/scs/v2"
	storage "github.com/alienix2/sensor_info/pkg/storage/central_database"
	"github.com/alienix2/sensor_visualization/cmd/web/mocks"
	"github.com/alienix2/sensor_visualization/internal/models"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type ACL struct {
	RW     *bool  `gorm:"type:tinyint(1)"`
	Topic  string `gorm:"type:varchar(255)"`
	ID     uint   `gorm:"primaryKey"`
	UserID int    `gorm:"index"`
}

func (ACL) TableName() string {
	return "acl"
}

func TestRoutes(t *testing.T) {
	ctx := context.Background()
	dsn, err := storage.SetupContainer(ctx)
	assert.NoError(t, err, "Failed to set up the database container")

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	assert.NoError(t, err, "Failed to connect to the database")

	err = db.AutoMigrate(&models.Topic{}, &models.Account{}, &ACL{}, &storage.MessageData{})
	assert.NoError(t, err, "Failed to migrate database schema")

	staticDir := "./ui/static"
	err = os.MkdirAll(staticDir+"/css", os.ModePerm)
	assert.NoError(t, err, "Failed to create fake static directory")

	defer os.RemoveAll("./ui")

	err = os.WriteFile(staticDir+"/css/main.css", []byte("body { background: #fff; }"), os.ModePerm)
	assert.NoError(t, err, "Failed to create fake static file")

	app := &application{
		sessionManager: scs.New(),
		topics:         &models.TopicModel{DB: db},
		accounts:       &models.AccountModel{DB: db},
		templateCache:  mocks.FakeTemplateCache(),
		infoLog:        log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime),
		errorLog:       log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile),
	}
	defer storage.TerminateMariaDBContainer(ctx)
	ts := httptest.NewTLSServer(app.routes())
	defer ts.Close()

	tests := []struct {
		name       string
		url        string
		method     string
		wantStatus int
	}{
		{"Home", "/", "GET", http.StatusOK},
		{"MessagesByTopic", "/messages/view?topic=faketopic", "GET", http.StatusOK},
		{"ControlMessageCreate", "/messages/sendcommand", "GET", http.StatusOK},
		{"ControlMessageCreatePost", "/messages/sendcommand", "POST", http.StatusOK},
		{"UserLogin", "/user/login", "GET", http.StatusOK},
		{"UserLoginPost", "/user/login", "POST", http.StatusOK},
		{"UserLogout", "/user/logout", "GET", http.StatusOK},
		{"StaticFiles", "/static/css/main.css", "GET", http.StatusOK},
		{"StaticFiles", "/static/", "GET", http.StatusNotFound},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest(tt.method, ts.URL+tt.url, nil)
			assert.NoError(t, err, "Failed to create HTTP request")

			resp, err := ts.Client().Do(req)
			assert.NoError(t, err, "Failed to make HTTP request")
			defer resp.Body.Close()

			assert.Equal(t, tt.wantStatus, resp.StatusCode, "Unexpected status code for "+tt.name)
		})
	}
}
