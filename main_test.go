package main

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"todox/database"
	"todox/handlers"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

func TestTodoAPI(t *testing.T) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Set up a test router
	r := setupTestRouter()

	// Create a new todo
	createTodoTest(t, r)

	// List all todos
	listTodosTest(t, r)

	// Update a todo
	updateTodoTest(t, r)

	// Remove a todo
	deleteTodoTest(t, r)

	database.DB.Exec(context.Background(), "TRUNCATE TABLE todos RESTART IDENTITY;")
}

func setupTestRouter() *gin.Engine {
	// Initialize the database for testing
	database.Init()

	// Create a new Gin router
	r := gin.Default()

	// Define API routes
	r.POST("/todos", handlers.CreateTodo)
	r.GET("/todos", handlers.ListTodos)
	r.PUT("/todos/:id", handlers.UpdateTodo)
	r.DELETE("/todos/:id", handlers.DeleteTodo)

	return r
}

func createTodoTest(t *testing.T, r *gin.Engine) {
	// Define the test cases
	tests := []struct {
		name           string
		requestBody    string
		expectedStatus int
	}{
		{
			name:           "Valid Todo",
			requestBody:    `{"title": "Test Title", "description": "Test Description", "status": true}`,
			expectedStatus: http.StatusCreated,
		},
		{
			name:           "Valid Todo",
			requestBody:    `{"title": "Test Title", "description": "Test Description", "status": true}`,
			expectedStatus: http.StatusCreated,
		},
		{
			name:           "Missing Title",
			requestBody:    `{"description": "Test Description", "status": true}`,
			expectedStatus: http.StatusBadRequest,
		},
		// Add more test cases as needed
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a request
			req, err := http.NewRequest("POST", "/todos", bytes.NewBufferString(tt.requestBody))
			if err != nil {
				t.Fatal(err)
			}

			// Set the request header to JSON
			req.Header.Set("Content-Type", "application/json")

			// Create a response recorder
			rec := httptest.NewRecorder()

			// Serve the request to the router
			r.ServeHTTP(rec, req)

			// Check the status code
			assert.Equal(t, tt.expectedStatus, rec.Code)

			// If the response status is OK, validate the response body
			if tt.expectedStatus == http.StatusCreated {
				var response handlers.Response
				err := json.Unmarshal(rec.Body.Bytes(), &response)
				assert.Nil(t, err)

				// Add your assertions for the response body here
				assert.NotZero(t, response.Todo.ID)
				assert.Equal(t, "Test Title", response.Todo.Title)
				assert.Equal(t, "Test Description", response.Todo.Description)
				assert.Equal(t, true, response.Todo.Status)
				assert.NotEmpty(t, response.Todo.CreatedDate)
			}
		})
	}
}

func updateTodoTest(t *testing.T, r *gin.Engine) {
	// Define the test cases
	tests := []struct {
		name           string
		todoID         string
		requestBody    string
		expectedStatus int
	}{
		{
			name:           "Valid Todo Update",
			todoID:         "1",
			requestBody:    `{"title": "Updated Title", "description": "Updated Description", "status": true}`,
			expectedStatus: http.StatusCreated,
		},
		{
			name:           "Invalid Todo ID",
			todoID:         "invalid",
			requestBody:    `{"title": "Updated Title", "description": "Updated Description", "status": true}`,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Todo Not Found",
			todoID:         "999",
			requestBody:    `{"title": "Updated Title", "description": "Updated Description", "status": true}`,
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "Invalid JSON",
			todoID:         "1",
			requestBody:    `{"invalid": "data"`,
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a request
			req, err := http.NewRequest("PUT", "/todos/"+tt.todoID, bytes.NewBufferString(tt.requestBody))
			if err != nil {
				t.Fatal(err)
			}

			// Set the request header to JSON
			req.Header.Set("Content-Type", "application/json")

			// Create a response recorder
			rec := httptest.NewRecorder()

			// Serve the request to the router
			r.ServeHTTP(rec, req)

			// Check the status code
			assert.Equal(t, tt.expectedStatus, rec.Code)

			// If the response status is OK, validate the response body
			if tt.expectedStatus == http.StatusCreated {
				var response handlers.Response
				err := json.Unmarshal(rec.Body.Bytes(), &response)
				assert.Nil(t, err)

				// Add your assertions for the response body here
				assert.Equal(t, tt.todoID, strconv.Itoa(response.Todo.ID))
				assert.Equal(t, "Updated Title", response.Todo.Title)
				assert.Equal(t, "Updated Description", response.Todo.Description)
				assert.Equal(t, true, response.Todo.Status)
				assert.NotEmpty(t, response.Todo.CreatedDate)
			}
		})
	}
}

func listTodosTest(t *testing.T, r *gin.Engine) {
	// Define the test cases
	tests := []struct {
		name           string
		queryParams    string
		expectedStatus int
		expectedLength int
	}{
		{
			name:           "List Todos",
			queryParams:    "",
			expectedStatus: http.StatusOK,
			expectedLength: 2, // Update with the expected number of todos in your test data
		},
		{
			name:           "List Todos Pagination",
			queryParams:    "pageSize=1&page=1",
			expectedStatus: http.StatusOK,
			expectedLength: 1, // Update with the expected number of todos in your test data
		},
		{
			name:           "List Todos Date",
			queryParams:    "startDate=2021-12-04T19:52:36Z",
			expectedStatus: http.StatusOK,
			expectedLength: 2, // Update with the expected number of todos in your test data
		},
		{
			name:           "List Todos Date",
			queryParams:    "endDate=2021-12-04T19:52:36Z",
			expectedStatus: http.StatusOK,
			expectedLength: 0, // Update with the expected number of todos in your test data
		},
		{
			name:           "List Todos Status",
			queryParams:    "status=true",
			expectedStatus: http.StatusOK,
			expectedLength: 2, // Update with the expected number of todos in your test data
		},
		{
			name:           "List Todos Status",
			queryParams:    "status=false",
			expectedStatus: http.StatusOK,
			expectedLength: 0, // Update with the expected number of todos in your test data
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a request
			req, err := http.NewRequest("GET", "/todos?"+tt.queryParams, nil)
			if err != nil {
				t.Fatal(err)
			}

			// Create a response recorder
			rec := httptest.NewRecorder()

			// Serve the request to the router
			r.ServeHTTP(rec, req)

			// Check the status code
			assert.Equal(t, tt.expectedStatus, rec.Code)

			// If the response status is OK, validate the response body
			if tt.expectedStatus == http.StatusOK {
				var response handlers.TodosResponse
				err := json.Unmarshal(rec.Body.Bytes(), &response)
				assert.Nil(t, err)

				// Add your assertions for the response body here
				assert.Len(t, response.Todos, tt.expectedLength)
			}
		})
	}
}

func deleteTodoTest(t *testing.T, r *gin.Engine) {
	tests := []struct {
		name           string
		todoID         string
		expectedStatus int
	}{
		{
			name:           "Valid Todo Deletion",
			todoID:         "1",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Invalid Todo ID",
			todoID:         "invalid",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Todo Not Found",
			todoID:         "999",
			expectedStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a request
			req, err := http.NewRequest("DELETE", "/todos/"+tt.todoID, nil)
			if err != nil {
				t.Fatal(err)
			}

			// Create a response recorder
			rec := httptest.NewRecorder()

			// Serve the request to the router
			r.ServeHTTP(rec, req)

			// Check the status code
			assert.Equal(t, tt.expectedStatus, rec.Code)

			// If the response status is OK, validate the response body
			if tt.expectedStatus == http.StatusOK {
				var response map[string]string
				err := json.Unmarshal(rec.Body.Bytes(), &response)
				assert.Nil(t, err)

				// Add your assertions for the response body here
				assert.Equal(t, "Todo deleted successfully", response["message"])
			}
		})
	}
}
