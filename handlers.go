package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func TodoCreate(w http.ResponseWriter, r *http.Request) {
	var todo Todo
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		log.Fatalf("Error reading from request: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := r.Body.Close(); err != nil {
		log.Fatalf("Error when closing body of request: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := json.Unmarshal(body, &todo); err != nil {
		log.Fatalf("Error when unmarshaling json into struct: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	t := RepoCreateTodo(todo)
	w.Header().Set("Content-Type", "application/json;charset=UTF-8")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(t); err != nil {
		log.Fatalf("Error when encoding json")
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func TodoIndex(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(todos); err != nil {
		log.Fatalf("Error when encoding json")
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func Index(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	fmt.Fprintln(w, "Welcome!")
}

func TodoShow(w http.ResponseWriter, r *http.Request) {
	todoId := mux.Vars(r)["todoId"]
	w.Header().Set("Content-Type", "text/html; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	id, err := strconv.Atoi(todoId)
	if err != nil {
		log.Fatalf("Error while converting string to int %s", err)
		w.WriteHeader(http.StatusInternalServerError)
	}

	for _, todo := range todos {
		if todo.Id == id {
			if err := json.NewEncoder(w).Encode(todo); err != nil {
				log.Fatalf("Error when encoding json")
				w.WriteHeader(http.StatusInternalServerError)
			}
			return
		}
	}

	log.Fatalf("Id not in use %d", id)
	w.WriteHeader(http.StatusInternalServerError)
}
