package main

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// Game model
type Game struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
	Genre string `json:"genre"`
	Price int    `json:"price"`
}

// In-memory database (slice of games)
// TODO: use a database (sqlite, maybe)
var games = []Game{
	{ID: 1, Title: "Legend of the Dragoon", Genre: "RPG", Price: 50},
	{ID: 2, Title: "Skyrim", Genre: "RPG", Price: 60},
	{ID: 3, Title: "Valkyria Chronicles", Genre: "Strategy RPG", Price: 40},
}

// GetGames retrieves all games
func GetGames(c *gin.Context) {
	c.JSON(http.StatusOK, games)
}

// CreateGame adds a new game
func CreateGame(c *gin.Context) {
	var newGame Game
	if err := c.ShouldBindJSON(&newGame); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	games = append(games, newGame)
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

	// Search for the game by ID
	for _, game := range games {
		if game.ID == gameID {
			c.JSON(http.StatusOK, game)
			return
		}
	}

	c.JSON(http.StatusNotFound, gin.H{"message": "Game not found"})
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

	// Search for the game by ID and update it
	for i, game := range games {
		if game.ID == gameID {
			games[i] = updatedGame
			games[i].ID = gameID // Preserve the original ID
			c.JSON(http.StatusOK, updatedGame)
			return
		}
	}

	c.JSON(http.StatusNotFound, gin.H{"message": "Game not found"})
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

	// Search for the game by ID and delete it
	for i, game := range games {
		if game.ID == gameID {
			games = append(games[:i], games[i+1:]...)
			c.JSON(http.StatusOK, gin.H{"message": "Game deleted"})
			return
		}
	}

	c.JSON(http.StatusNotFound, gin.H{"message": "Game not found"})
}

func main() {
	// Initialize Gin router - Logging included!
	r := gin.Default()

	// Define routes
	r.GET("/games", GetGames)
	r.POST("/games", CreateGame)
	r.GET("/games/:id", GetGameByID)
	r.PUT("/games/:id", UpdateGame)
	r.DELETE("/games/:id", DeleteGame)

	// Start the server
	r.Run(":8080") // Runs on http://localhost:8080
}
