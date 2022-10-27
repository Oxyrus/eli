package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/Oxyrus/eli/internal/models"
	"github.com/gorilla/mux"
	"mime"
	"net/http"
	"strconv"
	"time"
)

type TaskServer struct {
	store *models.TaskStore
}

func NewTaskServer() *TaskServer {
	store := models.New()
	return &TaskServer{store}
}

func (ts *TaskServer) CreateTaskHandler(w http.ResponseWriter, r *http.Request) {
	type CreateTaskRequest struct {
		Text string    `json:"text"`
		Tags []string  `json:"tags"`
		Due  time.Time `json:"due"`
	}

	type CreateTaskResponse struct {
		Id int `json:"id"`
	}

	contentType := r.Header.Get("Content-Type")
	mediaType, _, err := mime.ParseMediaType(contentType)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if mediaType != "application/json" {
		http.Error(w, "expect application/json Content-Type", http.StatusUnsupportedMediaType)
		return
	}

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	var rt CreateTaskRequest
	if err := dec.Decode(&rt); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	id := ts.store.CreateTask(rt.Text, rt.Tags, rt.Due)
	renderJSON(w, CreateTaskResponse{Id: id})
}

func (ts *TaskServer) GetAllTasksHandler(w http.ResponseWriter, r *http.Request) {
	allTasks := ts.store.GetAllTasks()
	renderJSON(w, allTasks)
}

func (ts *TaskServer) GetTaskHandler(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	task, err := ts.store.GetTask(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNoContent)
		return
	}
	renderJSON(w, task)
}

func (ts *TaskServer) DeleteAllTasksHandler(w http.ResponseWriter, r *http.Request) {
	ts.store.DeleteAllTasks()
	w.WriteHeader(http.StatusNoContent)
}

func (ts *TaskServer) DeleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	err := ts.store.DeleteTask(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
	}
}

func (ts *TaskServer) GetByTagHandler(w http.ResponseWriter, r *http.Request) {
	tag := mux.Vars(r)["tag"]
	tasks := ts.store.GetTasksByTag(tag)
	renderJSON(w, tasks)
}

func (ts *TaskServer) GetByDueDateHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	badRequestError := func() {
		http.Error(w, fmt.Sprintf("expected /due/<year>/<month>/<day>/, got %v", r.URL.Path), http.StatusBadRequest)
	}

	year, _ := strconv.Atoi(vars["year"])
	month, _ := strconv.Atoi(vars["month"])
	if month < int(time.January) || month > int(time.December) {
		badRequestError()
		return
	}

	day, _ := strconv.Atoi(vars["day"])
	tasks := ts.store.GetTasksByDueDate(year, time.Month(month), day)
	renderJSON(w, tasks)
}

func renderJSON(w http.ResponseWriter, v interface{}) {
	js, err := json.Marshal(v)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}
