package main

import (
	"github.com/Oxyrus/eli/internal/server"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func main() {
	router := mux.NewRouter()
	router.StrictSlash(true)
	sv := server.NewTaskServer()

	router.HandleFunc("/task/", sv.CreateTaskHandler).Methods("POST")

	err := http.ListenAndServe(":3000", router)
	if err != nil {
		log.Fatal(err)
	}
}
