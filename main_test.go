package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/stretchr/testify/assert"
)

// Helper function to initialize a fresh SQLite DB for testing
func setupTestDB() (*gorm.DB, error) {
	// Create an in-memory SQLite database
	db, err := gorm.Open("sqlite3", ":memory:") // Use ":memory:" for in-memory DB in tests
	if err != nil {
		return nil, err
	}

	// Migrate schema
	db.AutoMigrate(&Game{})
	return db, nil
}

// TestCreateGame tests the creation of a new game
func TestCreateGame(t *testing.T) {
	// Setup the test database
	db, err := setupTestDB()
	if err != nil {
		t.Fatalf("Failed to connect to the database: %v", err)
	}
	defer db.Close()

	// Initialize Gin
	gin.SetMode(gin.TestMode)
	r := gin.Default()

	// Inject the DB into the context
	r.POST("/games", func(c *gin.Context) {
		var newGame Game
		if err := c.ShouldBindJSON(&newGame); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		// Use GORM to save the game
		if err := db.Create(&newGame).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create game"})
			return
		}
		c.JSON(http.StatusCreated, newGame)
	})

	// Send a POST request to create a game
	gameJSON := `{"title": "The Witcher 3", "genre": "RPG", "price": 50}`
	req, _ := http.NewRequest(http.MethodPost, "/games", bytes.NewBufferString(gameJSON))
	req.Header.Set("Content-Type", "application/json")

	// Create a response recorder
	w := httptest.NewRecorder()

	// Perform the request
	r.ServeHTTP(w, req)

	// Check that the game was created
	assert.Equal(t, http.StatusCreated, w.Code)

	// Check the game is in the database
	var game Game
	if err := db.First(&game).Error; err != nil {
		t.Fatalf("Failed to find game in the database: %v", err)
	}

	assert.Equal(t, "The Witcher 3", game.Title)
	assert.Equal(t, "RPG", game.Genre)
	assert.Equal(t, 50, game.Price)
}

// TestGetGameByID tests retrieving a game by its ID
func TestGetGameByID(t *testing.T) {
	// Setup the test database
	db, err := setupTestDB()
	if err != nil {
		t.Fatalf("Failed to connect to the database: %v", err)
	}
	defer db.Close()

	// Seed the database with a game
	game := Game{Title: "Skyrim", Genre: "RPG", Price: 60}
	db.Create(&game)

	// Initialize Gin
	gin.SetMode(gin.TestMode)
	r := gin.Default()

	// Inject the DB into the context
	r.GET("/games/:id", func(c *gin.Context) {
		id := c.Param("id")
		var game Game
		if err := db.First(&game, id).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"message": "Game not found"})
			return
		}
		c.JSON(http.StatusOK, game)
	})

	// Send a GET request to retrieve the game
	req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/games/%d", game.ID), nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// Check that the response code is OK
	assert.Equal(t, http.StatusOK, w.Code)

	// Check the content of the response
	var retrievedGame Game
	err = json.Unmarshal(w.Body.Bytes(), &retrievedGame)
	if err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	assert.Equal(t, game.Title, retrievedGame.Title)
	assert.Equal(t, game.Genre, retrievedGame.Genre)
	assert.Equal(t, game.Price, retrievedGame.Price)
}
