package handlers

import (
	"context"
	"net/http"
	"strconv"

	"todox/database"

	"github.com/gin-gonic/gin"
)

// @Summary Delete a todo
// @Description Remove a todo by ID
// @Tags todos
// @Accept json
// @Produce json
// @Param id path int true "Todo ID"
// @Success 200 {object} SuccessResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /todos/{id} [delete]
func DeleteTodo(c *gin.Context) {
	todoID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid todo ID"})
		return
	}

	_, err = database.DB.Exec(context.Background(), "DELETE FROM todos WHERE id = $1", todoID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete todo"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Todo deleted successfully"})
}
