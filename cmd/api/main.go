package main

import (
	"github.com/Oxyrus/eli/internal/handlers"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func main() {
	router := mux.NewRouter()
	router.StrictSlash(true)
	sv := handlers.NewTaskServer()

	router.HandleFunc("/tasks/", sv.CreateTaskHandler).Methods("POST")
	router.HandleFunc("/tasks/", sv.GetAllTasksHandler).Methods("GET")
	router.HandleFunc("/tasks/", sv.DeleteAllTasksHandler).Methods("DELETE")
	router.HandleFunc("/tasks/{id:[0-9]+}/", sv.GetTaskHandler).Methods("GET")
	router.HandleFunc("/tasks/{id:[0-9]+}/", sv.DeleteTaskHandler).Methods("DELETE")
	router.HandleFunc("/tasks/due/{year:[0-9]+}/{month:[0-9]+}/{day:[0-9]+}/", sv.GetByDueDateHandler).Methods("GET")

	err := http.ListenAndServe(":3000", router)
	if err != nil {
		log.Fatal(err)
	}
}
