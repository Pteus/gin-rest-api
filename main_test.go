package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func setUpRouter() *gin.Engine {
	// Setup the router with routes
	r := gin.Default()
	r.GET("/games", GetGames)
	r.POST("/games", CreateGame)
	r.GET("/games/:id", GetGameByID)
	r.PUT("/games/:id", UpdateGame)
	r.DELETE("/games/:id", DeleteGame)
	return r
}

// Test for GET /games
func TestGetGames(t *testing.T) {
	r := setUpRouter()

	// Perform GET request
	req, _ := http.NewRequest("GET", "/games", nil)
	w := performRequest(r, req)

	// Assert the status code and the response body
	assert.Equal(t, http.StatusOK, w.Code)

	// Assert response is a JSON array
	var games []Game
	err := json.Unmarshal(w.Body.Bytes(), &games)
	assert.Nil(t, err)
	assert.NotEmpty(t, games)
}

// Test for POST /games
func TestCreateGame(t *testing.T) {
	r := setUpRouter()

	// Game to create
	newGame := Game{
		Title: "The Witcher 3",
		Genre: "RPG",
		Price: 50,
	}

	// Convert to JSON
	gameJSON, _ := json.Marshal(newGame)

	// Perform POST request
	req, _ := http.NewRequest("POST", "/games", bytes.NewBuffer(gameJSON))
	req.Header.Set("Content-Type", "application/json")
	w := performRequest(r, req)

	// Assert the status code and the response body
	assert.Equal(t, http.StatusCreated, w.Code)

	// Check that the response has the new game data
	var createdGame Game
	err := json.Unmarshal(w.Body.Bytes(), &createdGame)
	assert.Nil(t, err)
	assert.Equal(t, newGame.Title, createdGame.Title)
	assert.Equal(t, newGame.Genre, createdGame.Genre)
	assert.Equal(t, newGame.Price, createdGame.Price)
}

// Test for GET /games/:id
func TestGetGameByID(t *testing.T) {
	r := setUpRouter()

	// Perform GET request for the game with ID 1
	req, _ := http.NewRequest("GET", "/games/1", nil)
	w := performRequest(r, req)

	// Assert the status code and the response body
	assert.Equal(t, http.StatusOK, w.Code)

	// Check that the returned game has ID 1
	var game Game
	err := json.Unmarshal(w.Body.Bytes(), &game)
	assert.Nil(t, err)
	assert.Equal(t, 1, game.ID)
}

// Test for PUT /games/:id
func TestUpdateGame(t *testing.T) {
	r := setUpRouter()

	// Game to update (ID 1)
	updatedGame := Game{
		Title: "The Witcher 3 - Enhanced Edition",
		Genre: "RPG",
		Price: 55,
	}

	// Convert to JSON
	gameJSON, _ := json.Marshal(updatedGame)

	// Perform PUT request for game with ID 1
	req, _ := http.NewRequest("PUT", "/games/1", bytes.NewBuffer(gameJSON))
	req.Header.Set("Content-Type", "application/json")
	w := performRequest(r, req)

	// Assert the status code and the response body
	assert.Equal(t, http.StatusOK, w.Code)

	// Check that the returned game has the updated data
	var game Game
	err := json.Unmarshal(w.Body.Bytes(), &game)
	assert.Nil(t, err)
	assert.Equal(t, updatedGame.Title, game.Title)
	assert.Equal(t, updatedGame.Genre, game.Genre)
	assert.Equal(t, updatedGame.Price, game.Price)
}

// Test for DELETE /games/:id
func TestDeleteGame(t *testing.T) {
	r := setUpRouter()

	// Perform DELETE request for the game with ID 1
	req, _ := http.NewRequest("DELETE", "/games/1", nil)
	w := performRequest(r, req)

	// Assert the status code and the response body
	assert.Equal(t, http.StatusOK, w.Code)

	// Perform GET request to verify that the game is deleted
	req, _ = http.NewRequest("GET", "/games/1", nil)
	w = performRequest(r, req)

	// Assert that the game is not found
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func performRequest(r *gin.Engine, req *http.Request) *httptest.ResponseRecorder {
	// Perform the request using Gin's test mode
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}
