package main

import (
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/alienix2/sensor_visualization/internal/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type application struct {
	errorLog      *log.Logger
	infoLog       *log.Logger
	topics        *models.TopicModel
	accounts      *models.AccountModel
	templateCache map[string]*template.Template
}

func main() {
	addr := flag.String("addr", ":4000", "HTTP network address")
	dsn := flag.String("database_path", "mqtt_admin:Panzerotto@tcp(localhost:3306)/mqtt_users?parseTime=true", "Path to the MySQL database")
	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	db, err := gorm.Open(mysql.Open(*dsn), &gorm.Config{})

	templateCache, err := newTemplateCache()

	fmt.Println(templateCache)
	app := application{
		errorLog:      errorLog,
		infoLog:       infoLog,
		topics:        &models.TopicModel{DB: db},
		accounts:      &models.AccountModel{DB: db},
		templateCache: templateCache,
	}

	srv := &http.Server{
		Addr:     *addr,
		ErrorLog: errorLog,
		Handler:  app.routes(),
	}

	log.Printf("Starting server on %s", *addr)
	err = srv.ListenAndServe()
	log.Fatal(err)
}
