package gas

import (
	"io"
	"os"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	easy "github.com/t-tomalak/logrus-easy-formatter"
)

func LogrusInit() {
	config := viper.New()
	config.SetConfigFile("logrus.env")
	config.AddConfigPath(".")

	err := config.ReadInConfig()
	if err != nil {
		logrus.Error("Error reading config file: ", err)
		return
	}

	// Get the log level from the configuration file or use a default
	logrusLevel := config.GetString("SET_LEVEL")
	if logrusLevel == "" {
		logrusLevel = "InfoLevel" // Set a default log level
	}

	// Parse the log level string into a logrus level
	level, err := logrus.ParseLevel(logrusLevel)
	if err != nil {
		logrus.Warn("Error parsing log level:", err)
		return
	}

	// Create a filename based on the current date
	logFileName := "log_" + time.Now().Format("2006-01-02") + ".txt"

	// Check if the file already exists
	if _, err := os.Stat(logFileName); err == nil {
		// If the file exists, append to it
		logrus.Info("Log file already exists. Appending to", logFileName)
	} else {
		// If the file does not exist, create it
		logrus.Info("Creating new log file:", logFileName)
	}

	// Open the file for writing or appending
	f, err := os.OpenFile(logFileName, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		logrus.Warn("Failed to create/open logfile:", logFileName)
		logrus.Warn("OpenFile function error:", err)
	}
	defer f.Close()

	// Set up the logrus logger
	logrus.SetOutput(io.MultiWriter(f, os.Stdout))
	logrus.SetLevel(level)
	logrus.SetFormatter(&easy.Formatter{
		TimestampFormat: "2006-01-02 15:04:05",
		LogFormat:       "[%lvl%]: %time% - %msg%\n",
	})
}
