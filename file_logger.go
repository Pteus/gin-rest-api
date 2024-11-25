package main

import (
	"os"

	"github.com/sirupsen/logrus"
)

// Logger is the global logger instance
var Logger = logrus.New()

// SetupFileLogger configures logging to a file
func SetupFileLogger() {
	// Create or append to the log file
	logFile, err := os.OpenFile("gin.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		Logger.Fatal("Failed to open log file: ", err)
	}

	// Set log output to the file
	Logger.SetOutput(logFile)

	// Set log format and log level
	Logger.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true, // Include timestamps
	})
	Logger.SetLevel(logrus.InfoLevel) // Adjust the log level as needed
}
