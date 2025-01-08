package main

import (
	"net/http"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	userID := app.sessionManager.GetInt(r.Context(), "authenticatedUserID")
	topics, err := app.topics.GetTopics(userID)
	if err != nil {
		app.errorLog.Println(err.Error())
	}

	data := app.newTemplateData(r)
	data.Topics = topics
	app.render(w, http.StatusOK, "home.tmpl", data)
}

func (app *application) messagesByTopic(w http.ResponseWriter, r *http.Request) {
	topic := r.URL.Query().Get("topic")
	if topic == "" {
		app.errorLog.Println("Topic is required")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	userID := app.sessionManager.GetInt(r.Context(), "authenticatedUserID")
	messages, err := app.topics.GetMessagesByTopic(userID, topic)
	if err != nil {
		app.errorLog.Println(err.Error())
		http.Error(w, "Unable to fetch messages", http.StatusInternalServerError)
		return
	}

	data := app.newTemplateData(r)
	data.Topic = topic
	data.MessageData = messages
	app.render(w, http.StatusOK, "message_data.tmpl", data)
}
