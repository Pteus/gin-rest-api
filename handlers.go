package main

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// GetGames retrieves all games
func GetGames(c *gin.Context) {
	var games []Game
	if err := db.Find(&games).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch games"})
		return
	}
	c.JSON(http.StatusOK, games)
}

// CreateGame adds a new game
func CreateGame(c *gin.Context) {
	var newGame Game
	if err := c.ShouldBindJSON(&newGame); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := db.Create(&newGame).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create game"})
		return
	}

	c.JSON(http.StatusCreated, newGame)
}

// GetGameByID retrieves a game by its ID
func GetGameByID(c *gin.Context) {
	id := c.Param("id")

	// Convert the id to an int
	gameID, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	var game Game
	if err := db.First(&game, gameID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Game not found"})
		return
	}

	c.JSON(http.StatusOK, game)
}

// UpdateGame updates an existing game by its ID
func UpdateGame(c *gin.Context) {
	id := c.Param("id")
	var updatedGame Game

	// Convert the id to an int
	gameID, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	// Bind the JSON request body to the Game struct
	if err := c.ShouldBindJSON(&updatedGame); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Find the game by ID and update it
	var game Game
	if err := db.First(&game, gameID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Game not found"})
		return
	}

	game.Title = updatedGame.Title
	game.Genre = updatedGame.Genre
	game.Price = updatedGame.Price

	// Save the updated game back to the database
	if err := db.Save(&game).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update game"})
		return
	}

	c.JSON(http.StatusOK, game)
}

// DeleteGame deletes a game by its ID
func DeleteGame(c *gin.Context) {
	id := c.Param("id")

	// Convert the id to an int
	gameID, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	var game Game
	if err := db.First(&game, gameID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Game not found"})
		return
	}

	// Delete the game from the database
	if err := db.Delete(&game).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete game"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Game deleted"})
}

// Login is a mock route to simulate user login and return a JWT token
func Login(c *gin.Context) {
	var loginInfo struct {
		UserID int `json:"userID"`
	}

	// Bind the incoming JSON request to the struct
	if err := c.ShouldBindJSON(&loginInfo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Generate JWT token for the user
	token, err := GenerateJWT(loginInfo.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate token"})
		return
	}

	// Return the JWT token to the user
	c.JSON(http.StatusOK, gin.H{"token": token})
}
