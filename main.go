package main

import (
	"log"
	"os"

	"todox/database"
	"todox/handlers"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	docs "todox/docs"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Read the Gin app api port from environment variable or use a default value (8080)
	apiPort := os.Getenv("API_PORT")
	if apiPort == "" {
		apiPort = "8080"
	}

	// Initialize the database connection
	database.Init()

	// Create a new Gin router
	r := gin.Default()

	docs.SwaggerInfo.BasePath = "/"

	// Define API routes
	r.POST("/todos", handlers.CreateTodo)
	r.GET("/todos", handlers.ListTodos)
	r.PUT("/todos/:id", handlers.UpdateTodo)
	r.DELETE("/todos/:id", handlers.DeleteTodo)

	// Swagger route
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Run the server
	r.Run(":" + apiPort)
}
