package main

import (
	"encoding/json"
	"log"
	"mime"
	"net/http"
	"time"

	gorestserver "example.com/go-rest-server/stdlib"
)

type taskServer struct {
	store *gorestserver.TaskStore
}

func NewTaskServer() *taskServer {
	store := gorestserver.New()
	return &taskServer{store: store}
}

func renderJSON(w http.ResponseWriter, v interface{}) {
	js, err := json.Marshal(v)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func (ts *taskServer) createTaskHandler(w http.ResponseWriter, req *http.Request) {

	log.Printf("handling task create at %s\n", req.URL.Path)

	type RequestTask struct {
		Text string    `json:"text"`
		Tags []string  `json:"tags"`
		Due  time.Time `json:"due"`
	}

	type ResponseId struct {
		Id int `json:"id"`
	}

	// Enforce a JSON Content-Type.
	contentType := req.Header.Get("Content-Type")
	mediatype, _, err := mime.ParseMediaType(contentType)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if mediatype != "application/json" {
		http.Error(w, "expect application/json Content-Type", http.StatusUnsupportedMediaType)
		return
	}

	dec := json.NewDecoder(req.Body)
	dec.DisallowUnknownFields()
	var rt RequestTask
	if err := dec.Decode(&rt); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	id := ts.store.CreateTask(rt.Text, rt.Tags, rt.Due)
	renderJSON(w, ResponseId{Id: id})

}

func main() {
	// mux := http.NewServeMux()
	// server := NewTaskServer()

	// mux.HandleFunc("POST /task/", server.createTaskHandler)
	// mux.HandleFunc("GET /task/", server.getAllTasksHandler)
	// mux.HandleFunc("DELETE /task/", server.deleteAllTasksHandler)
	// mux.HandleFunc("GET /task/{id}/", server.getTaskHandler)
	// mux.HandleFunc("DELETE /task/{id}/", server.deleteTaskHandler)
	// mux.HandleFunc("GET /tag/{tag}/", server.tagHandler)
	// mux.HandleFunc("GET /due/{year}/{month}/{day}/", server.dueHandler)

	// log.Fatal(http.ListenAndServe("localhost:"+os.Getenv("SERVERPORT"), mux))

}
