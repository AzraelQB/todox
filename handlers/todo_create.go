package handlers

import (
	"context"
	"net/http"
	"time"

	"todox/database"

	"github.com/gin-gonic/gin"
)

type TodoRequest struct {
	Title       string `json:"title" validate:"required"`
	Description string `json:"description"`
	Status      bool   `json:"status"`
}

type TodoResponse struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Status      bool   `json:"status"`
	CreatedDate string `json:"created_date"`
}

// @Summary Create a new todo
// @Description Create a new todo with the provided data
// @Tags todos
// @Accept json
// @Produce json
// @Param input body TodoRequest true "Todo body"
// @Success 201 {object} Response
// @Failure 400 {object} ErrorResponse
// @Router /todos [post]
func CreateTodo(c *gin.Context) {
	var json TodoRequest

	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate the body struct using the validator package
	if err := validate.Struct(json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var id int
	var createdDate time.Time
	err := database.DB.QueryRow(context.Background(), `
        INSERT INTO todos (title, description, status) VALUES ($1, $2, $3) RETURNING id, created_date
    `, json.Title, json.Description, json.Status).Scan(&id, &createdDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create todo"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"todo": &TodoResponse{
		ID:          id,
		Title:       json.Title,
		Description: json.Description,
		Status:      json.Status,
		CreatedDate: createdDate.Format(time.RFC3339),
	}})
}
