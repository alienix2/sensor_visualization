package main

import (
	"crypto/tls"
	"flag"
	"html/template"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/alexedwards/scs/gormstore"
	"github.com/alexedwards/scs/v2"
	tls_config "github.com/alienix2/sensor_info/pkg/tls_config"
	"github.com/alienix2/sensor_visualization/internal/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type application struct {
	errorLog       *log.Logger
	infoLog        *log.Logger
	topics         models.TopicModelInterface
	accounts       models.AccountModelInterface
	templateCache  map[string]*template.Template
	sessionManager *scs.SessionManager
	tlsConfig      *tls.Config
	mqttClients    sync.Map
	mqttBroker     string
}

func main() {
	addr := flag.String("addr", ":4000", "HTTP network address")
	dsn := flag.String("dsn", "mqtt_admin:Panzerotto@tcp(localhost:3306)/mqtt_users?parseTime=true", "Path to the MySQL database")
	broker := flag.String("broker", "tls://localhost:8883", "MQTT broker address")

	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	db, err := gorm.Open(mysql.Open(*dsn), &gorm.Config{})
	if err != nil {
		errorLog.Fatal(err)
	}

	templateCache, err := newTemplateCache()
	if err != nil {
		errorLog.Fatal(err)
	}

	sessionManager := scs.New()
	sessionManager.Store, err = gormstore.New(db)
	if err != nil {
		errorLog.Fatal(err)
	}
	sessionManager.Lifetime = 12 * time.Hour
	sessionManager.Cookie.Secure = true

	tlsConfig, err := tls_config.LoadCertificates("certifications/publisher.crt", "certifications/publisher.key", "certifications/ca.crt")
	if err != nil {
		errorLog.Fatal(err)
	}

	app := application{
		errorLog:       errorLog,
		infoLog:        infoLog,
		topics:         &models.TopicModel{DB: db},
		accounts:       &models.AccountModel{DB: db},
		templateCache:  templateCache,
		sessionManager: sessionManager,
		mqttBroker:     *broker,
		tlsConfig:      tlsConfig,
	}

	srv := &http.Server{
		Addr:     *addr,
		ErrorLog: errorLog,
		Handler:  app.routes(),
	}

	infoLog.Printf("Starting server on %s", *addr)
	err = srv.ListenAndServeTLS("certifications/publisher.crt", "certifications/publisher.key")
	errorLog.Fatal(err)
}
