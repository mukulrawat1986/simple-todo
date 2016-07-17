package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

func Index(w http.ResponseWriter, r *http.Request) {
	todos := Todos{
		Todo{Name: "Write Presentation"},
		Todo{Name: "Host Meetup"},
	}

	if err := json.NewEncoder(w).Encode(todos); err != nil {
		log.Fatalf("Error while encoding todo into json")
		os.Exit(1)
	}
}

func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Welcome!")
}

func TodoShow(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	todoId := vars("todoId")
	fmt.Fprintln(w, "Todo show: ", todoId)
}
