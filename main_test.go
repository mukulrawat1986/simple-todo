package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/mux"
)

func assertEqual(t *testing.T, a interface{}, b interface{}, message string) {
	if a == b {
		return
	}
	if len(message) == 0 {
		message = fmt.Sprintf("%v != %v", a, b)
	}
	t.Fatal(message)
}

// Test TodoCreate using ResponseRecorder
func TestTodoCreate(t *testing.T) {

	// create an instance of router
	router := mux.NewRouter()

	// register the handler
	router.
		Methods("Post").
		Path("/todos").
		Name("TodoCreate").
		Handler(http.HandlerFunc(TodoCreate))

	// create an HTTP Request
	todoJson := `{"name": "New Todo"}`

	req, err := http.NewRequest(
		"POST",
		"/todos",
		strings.NewReader(todoJson),
	)

	assertEqual(t, nil, err, "")

	// create a ResponseRecorder object
	w := httptest.NewRecorder()

	// Send the ResponseRecorder object and Request object to the
	// multiplexer
	router.ServeHTTP(w, req)

	// Inspect the ResponseRecorder object
	assertEqual(t, 201, w.Code, fmt.Sprintf("HTTP Status expected:201, got: %d", w.Code))

	// check the body of response
	var todo Todo
	err = json.NewDecoder(w.Body).Decode(&todo)
	assertEqual(t, nil, err, "The error is not nil")
	assertEqual(t, "New Todo", todo.Name, "The json returned is not equal")
}

// Create end-to-end test for TodoCreate using httptest Server
func TestTodoCreateClient(t *testing.T) {

	// create a router instance
	router := mux.NewRouter()

	router.
		Methods("POST").
		Path("/todos").
		Name("TodoCreate").
		Handler(http.HandlerFunc(TodoCreate))

	// create an HTTP server
	server := httptest.NewServer(router)
	defer server.Close()

	// create a Request object
	todoUrl := fmt.Sprintf("%s/todos", server.URL)
	todoJson := `{"name": "New Todo"}`
	request, err := http.NewRequest("POST", todoUrl, strings.NewReader(todoJson))

	// Send an HTTP request to the server
	res, err := http.DefaultClient.Do(request)

	assertEqual(t, nil, err, "")

	// Inspect the Response object
	assertEqual(t, 201, res.StatusCode, "")

	var todo Todo

	body, err := ioutil.ReadAll(res.Body)

	assertEqual(t, nil, err, "")

	err = res.Body.Close()

	assertEqual(t, nil, err, "")

	err = json.Unmarshal(body, &todo)

	assertEqual(t, nil, err, "Error while decoding")
	assertEqual(t, "New Todo", todo.Name, "")
}

// Test the response of TodoIndex using ResponseRecorder
func TestTodoIndex(t *testing.T) {

	// initialized todos
	todos = Todos{
		Todo{
			Name: "New Todo",
		},
	}

	// create an instance of the router
	router := mux.NewRouter()

	router.
		Methods("GET").
		Path("/todos").
		Name("TodoIndex").
		Handler(http.HandlerFunc(TodoIndex))

	// create an HTTP Request
	request, err := http.NewRequest("GET", "/todos", nil)
	assertEqual(t, nil, err, "Error while making request")

	// create a ResponseRecorder
	w := httptest.NewRecorder()

	// send the ResponseRecorder and Request object to the ServeHTTP
	router.ServeHTTP(w, request)

	// Inspect the ResponseRecorder
	assertEqual(t, 200, w.Code, fmt.Sprintf("Expected code: 200, got: %d", w.Code))

	var todos []Todo
	err = json.NewDecoder(w.Body).Decode(&todos)
	assertEqual(t, nil, err, fmt.Sprintf("Error while decoding json %s", err))
	assertEqual(t, "New Todo", todos[0].Name, fmt.Sprintf("Expected json: New Todo, got: %s", todos[0].Name))
}

// Create end-to-end test for TodoIndex
func TestTodoIndexClient(t *testing.T) {

}

// Test the Response of Index using ResponseRecorder
func TestIndex(t *testing.T) {

	// create a new instance of the router
	router := mux.NewRouter()

	router.
		Methods("GET").
		Path("/").
		Name("Index").
		Handler(http.HandlerFunc(Index))

	// create an HTTP request
	req, err := http.NewRequest("GET", "/", nil)
	assertEqual(t, nil, err, "Error while creating Request")

	// Create a ResponseRecorder
	w := httptest.NewRecorder()

	// Send the ResponseRecorder and Request object to the ServeHTTP
	router.ServeHTTP(w, req)

	// Inspect the ResponseRecorder
	assertEqual(t, 200, w.Code, fmt.Sprintf("Expected code: 200, Got: %d", w.Code))
	assertEqual(t, "Welcome!", strings.TrimSpace(w.Body.String()), "")
}

// Test the Response of TodoShow using ResponseRecorder
func TestTodoShow(t *testing.T) {

	// Initialize the todos variable and currentid
	todos = Todos{}
	currentId = 0

	// Assign some new todo
	RepoCreateTodo(Todo{Name: "Hello"})
	RepoCreateTodo(Todo{Name: "World"})

	// Initialize a new instance of the router
	router := mux.NewRouter()

	router.
		Methods("GET").
		Path("/todos/{todoId}").
		Name("TodoShow").
		Handler(http.HandlerFunc(TodoShow))

	// create an HTTP Request
	req, err := http.NewRequest("GET", "/todos/1", nil)
	assertEqual(t, nil, err, "Error while making the request")

	// Create a ResponseRecorder
	w := httptest.NewRecorder()

	// Send the ResponseRecorder and Request object to ServeHTTP of router
	router.ServeHTTP(w, req)

	// Inspect the ResponseRecorder
	assertEqual(t, 200, w.Code, fmt.Sprintf("Expected code: 200, Got: %d", w.Code))

	var todo Todo

	err = json.NewDecoder(w.Body).Decode(&todo)
	assertEqual(t, nil, err, "Error when decoding json")
	assertEqual(t, "Hello", todo.Name, fmt.Sprintf("Expected string: Hello, got: %s", todo.Name))
}
