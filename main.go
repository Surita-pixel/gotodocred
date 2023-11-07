package main

import (
	"fmt"
	"log"
	"net/http"
	"github.com/gorilla/mux"
	"todoapi/services"
)

func indexRoute(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to my AaaaaPI")
}

func main() {
	router := mux.NewRouter().StrictSlash(true)

	router.HandleFunc("/", indexRoute)
	router.HandleFunc("/tasks", services.GetAllTasks).Methods("GET")
	router.HandleFunc("/tasks", services.CreateTask).Methods("POST")
	router.HandleFunc("/tasks/{id}", services.UpdateTask).Methods("PUT")
	router.HandleFunc("/tasks/{id}", services.GetOneTask).Methods("GET")
	router.HandleFunc("/tasks/{id}", services.DeleteTask).Methods("DELETE")


	log.Fatal(http.ListenAndServe(":5000", router))
}

