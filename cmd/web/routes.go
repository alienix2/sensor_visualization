package main

import (
	"net/http"
	"path/filepath"
)

func (app *application) routes() *http.ServeMux {
	mux := http.NewServeMux()

	fileServer := http.FileServer(neuteredFileSystem{http.Dir("./ui/static/")})
	mux.Handle("GET /static/", http.StripPrefix("/static", fileServer))

	sessionMux := http.NewServeMux()
	sessionMux.HandleFunc("GET /", app.home)
	sessionMux.HandleFunc("GET /messages/view", app.messagesByTopic)
	sessionMux.HandleFunc("GET /messages/sendcommand", app.ControlMessageCreate)
	sessionMux.HandleFunc("POST /messages/sendcommand", app.ControlMessageCreatePost)

	sessionMux.HandleFunc("GET /user/login", app.userLogin)
	sessionMux.HandleFunc("POST /user/login", app.userLoginPost)
	sessionMux.HandleFunc("GET /user/logout", app.userLogout)

	mux.Handle("/", app.sessionManager.LoadAndSave(sessionMux))

	return mux
}

type neuteredFileSystem struct {
	fs http.FileSystem
}

func (nfs neuteredFileSystem) Open(path string) (http.File, error) {
	f, err := nfs.fs.Open(path)
	if err != nil {
		return nil, err
	}

	s, err := f.Stat()
	if err != nil {
		return nil, err
	}

	if s.IsDir() {
		index := filepath.Join(path, "index.html")
		if _, err := nfs.fs.Open(index); err != nil {
			closeErr := f.Close()
			if closeErr != nil {
				return nil, closeErr
			}
			return nil, err
		}
	}

	return f, nil
}
