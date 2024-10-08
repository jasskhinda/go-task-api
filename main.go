package main

import (
	"encoding/json" // For working with JSON data
	"fmt"           // For formatted I/O
	"log"           // For logging errors
	"net/http"      // For HTTP server
	"strconv"       // For string conversions
	"strings"       // For string manipulation
)

type Task struct {
	ID          int    `json:"id"`          // Unique identifier
	Title       string `json:"title"`       // Task title
	Description string `json:"description"` // Task description
	Status      string `json:"status"`      // "pending" or "completed"
}

var tasks []Task // This will store all our tasks
var nextID = 1   // This will help us assign unique IDs to tasks

func createTask(w http.ResponseWriter, r *http.Request) {
	// Check if the method is POST
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	// Decode the JSON request body into a Task struct
	var newTask Task
	err := json.NewDecoder(r.Body).Decode(&newTask)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	// Assign a unique ID to the new task
	newTask.ID = nextID
	nextID++

	// Set the default status to "pending" if not provided
	if newTask.Status == "" {
		newTask.Status = "pending"
	}

	// Add the new task to the tasks slice
	tasks = append(tasks, newTask)

	// Return the created task as a response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated) // HTTP 201 Created
	json.NewEncoder(w).Encode(newTask)
}

func getTasks(w http.ResponseWriter, r *http.Request) {
	// Check if the method is GET
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	// Return all tasks as JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tasks)
}

func getTaskByID(w http.ResponseWriter, r *http.Request) {
	// Check if the method is GET
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract the ID from the URL
	idStr := strings.TrimPrefix(r.URL.Path, "/tasks/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid Task ID", http.StatusBadRequest)
		return
	}

	// Find the task with the matching ID
	for _, task := range tasks {
		if task.ID == id {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(task)
			return
		}
	}

	// If task not found, return an error
	http.Error(w, "Task Not Found", http.StatusNotFound)
}

func updateTask(w http.ResponseWriter, r *http.Request) {
	// Check if the method is PUT
	if r.Method != http.MethodPut {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract the ID from the URL
	idStr := strings.TrimPrefix(r.URL.Path, "/tasks/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid Task ID", http.StatusBadRequest)
		return
	}

	// Decode the JSON request body into a Task struct
	var updatedTask Task
	err = json.NewDecoder(r.Body).Decode(&updatedTask)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	// Find and update the task
	for i, task := range tasks {
		if task.ID == id {
			if updatedTask.Title != "" {
				tasks[i].Title = updatedTask.Title
			}
			if updatedTask.Description != "" {
				tasks[i].Description = updatedTask.Description
			}
			if updatedTask.Status != "" {
				tasks[i].Status = updatedTask.Status
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(tasks[i])
			return
		}
	}

	// If task not found, return an error
	http.Error(w, "Task Not Found", http.StatusNotFound)
}

func deleteTask(w http.ResponseWriter, r *http.Request) {
	// Check if the method is DELETE
	if r.Method != http.MethodDelete {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract the ID from the URL
	idStr := strings.TrimPrefix(r.URL.Path, "/tasks/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid Task ID", http.StatusBadRequest)
		return
	}

	// Find and delete the task
	for i, task := range tasks {
		if task.ID == id {
			// Remove the task from the slice
			tasks = append(tasks[:i], tasks[i+1:]...)
			w.WriteHeader(http.StatusNoContent) // HTTP 204 No Content
			return
		}
	}

	// If task not found, return an error
	http.Error(w, "Task Not Found", http.StatusNotFound)
}

func main() {
	// Route for "/tasks"
	http.HandleFunc("/tasks", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/tasks" {
			http.NotFound(w, r)
			return
		}

		switch r.Method {
		case http.MethodGet:
			getTasks(w, r)
		case http.MethodPost:
			createTask(w, r)
		default:
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})

	// Route for "/tasks/{id}"
	http.HandleFunc("/tasks/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			getTaskByID(w, r)
		case http.MethodPut:
			updateTask(w, r)
		case http.MethodDelete:
			deleteTask(w, r)
		default:
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})

	// Start the server on port 8080
	fmt.Println("Server is running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
