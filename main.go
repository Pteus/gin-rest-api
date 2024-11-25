package main

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

// DB instance
var db *gorm.DB
var err error

// Initialize and setup database connection
func initDB() {
	// Open a connection to SQLite
	db, err = gorm.Open("sqlite3", "./games.db")
	if err != nil {
		panic("failed to connect to the database")
	}

	// Migrate the schema (auto-create table based on the Game model)
	db.AutoMigrate(&Game{})
}

func main() {
	// Initialize Gin router - Logging included!
	r := gin.Default()

	// Initialize the database connection
	initDB()
	defer db.Close()

	// Initialize the file logger
	SetupFileLogger()
	r.Use(RequestLoggerMiddleware())

	// Public Routes (no JWT required)
	r.GET("/games", GetGames)
	r.POST("/games", CreateGame)

	// Route to login and get JWT token
	r.POST("/login", Login)

	// Authenticated Routes (JWT required)
	protected := r.Group("/protected") // Could make this /games and ajust the protected routes
	protected.Use(AuthMiddleware())    // Apply AuthMiddleware to this group

	protected.GET("/games/:id", GetGameByID)
	protected.PUT("/games/:id", UpdateGame)
	protected.DELETE("/games/:id", DeleteGame)

	// Start the server
	r.Run(":8080") // Runs on http://localhost:8080
}
