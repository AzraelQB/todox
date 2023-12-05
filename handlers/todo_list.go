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

// @Summary List todos
// @Description Get a list of todos
// @Tags todos
// @Accept json
// @Produce json
// @Param startDate query string false "Start date filter (RFC3339 format)"
// @Param endDate query string false "End date filter (RFC3339 format)"
// @Param status query bool false "Todo status filter"
// @Param pageSize query int false "Number of items per page"
// @Param page query int false "Page number"
// @Success 200 {array} TodosResponse
// @Failure 500 {object} ErrorResponse
// @Router /todos [get]
func ListTodos(c *gin.Context) {
	var todos []models.Todo

	query := `
        SELECT id, title, description, status, created_date FROM todos
        WHERE ($1::timestamp IS NULL OR created_date >= $1)
          AND ($2::timestamp IS NULL OR created_date <= $2)
          AND ($3::boolean IS NULL OR status = $3)
        ORDER BY created_date DESC
        LIMIT $4 OFFSET $5
    `

	startDateFilter, _ := time.Parse(time.RFC3339, c.Query("startDate"))
	endDateFilter, _ := time.Parse(time.RFC3339, c.DefaultQuery("endDate", time.Now().Format(time.RFC3339)))
	statusFilter := c.DefaultQuery("status", "")

	// Assuming you have a booleanValue variable based on the statusFilter
	var statusFilterBool bool
	var status *bool

	// Convert the string statusFilter to a boolean
	statusFilterBool, err := strconv.ParseBool(statusFilter)
	if err != nil {
		// Handle the error if the conversion fails
		status = nil
	} else {
		status = &statusFilterBool
	}

	pageSize, err := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid pageSize parameter"})
		return
	}

	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid page parameter"})
		return
	}

	offset := (page - 1) * pageSize

	rows, err := database.DB.Query(context.Background(), query, startDateFilter, endDateFilter, status, pageSize, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch todos"})
		return
	}
	defer rows.Close()

	for rows.Next() {
		var todo models.Todo
		if err := rows.Scan(&todo.ID, &todo.Title, &todo.Description, &todo.Status, &todo.CreatedDate); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch todos"})
			return
		}
		todos = append(todos, todo)
	}

	c.JSON(http.StatusOK, gin.H{"todos": todos})
}
