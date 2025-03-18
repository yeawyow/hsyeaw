package api

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres" // PostgreSQL driver
	"gorm.io/gorm"
)

// DB is the global variable for the database connection
var DB *gorm.DB

// SetupUserAPI sets up the routes for user-related API endpoints
func SetupUserAPI(router *gin.RouterGroup) {
	userAPI := router.Group("/user")
	{
		// Example route to test the API is working
		userAPI.GET("/", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "User API is working!"})
		})

		// Route to test the database connection
		userAPI.GET("/test", TestConnection)
	}
}

// TestConnection is a function to test the database connection
func TestConnection(c *gin.Context) {
	// Connection string (replace with your credentials)
	dsn := "host=192.168.0.90 user=hissync dbname=hissync password=Sync#6530 port=5432 sslmode=disable"
	// Open a connection to the database using the PostgreSQL driver
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		// If connection fails, send a JSON response with the error
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Error connecting to the database: %v", err)})
		return
	}

	// If successful, assign the connection to the global DB variable
	DB = db

	// Send a success response
	c.JSON(http.StatusOK, gin.H{"message": "Connected to PostgreSQL database successfully"})
}
