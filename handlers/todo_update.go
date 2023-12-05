package handlers

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"todox/database"
	"todox/models"

	"github.com/gin-gonic/gin"
)

// @Summary Update a todo
// @Description Update a todo by ID
// @Tags todos
// @Accept json
// @Produce json
// @Param id path int true "Todo ID"
// @Param input body TodoRequest true "Todo body"
// @Success 200 {object} Response
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /todos/{id} [put]
func UpdateTodo(c *gin.Context) {
	todoID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid todo ID"})
		return
	}

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

	// Check if the todo with the given ID exists
	var existingTodo models.Todo
	err = database.DB.QueryRow(context.Background(), "SELECT id FROM todos WHERE id = $1", todoID).Scan(&existingTodo.ID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Todo not found"})
		return
	}

	// Update the todo
	var createdDate time.Time
	err = database.DB.QueryRow(context.Background(), `
        UPDATE todos
        SET title = $1, description = $2, status = $3
        WHERE id = $4  RETURNING created_date
    `, json.Title, json.Description, json.Status, todoID).Scan(&createdDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update todo"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"todo": &TodoResponse{
		ID:          todoID,
		Title:       json.Title,
		Description: json.Description,
		Status:      json.Status,
		CreatedDate: createdDate.Format(time.RFC3339),
	}})
}
