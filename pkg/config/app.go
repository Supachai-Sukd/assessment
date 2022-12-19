package config

import (
	"os"
	"strconv"
	"time"
)

type App struct {
	Host        string
	Port        int
	Debug       bool
	ReadTimeout time.Duration
}

var app = &App{}

// AppCfg returns the default App configuration
func AppCfg() *App {
	return app
}

// LoadApp loads App configuration
func LoadApp() {
	app.Host = os.Getenv("APP_HOST")
	app.Port, _ = strconv.Atoi(os.Getenv("PORT"))
	app.Debug, _ = strconv.ParseBool(os.Getenv("APP_DEBUG"))
	timeOut, _ := strconv.Atoi(os.Getenv("APP_READ_TIMEOUT"))
	app.ReadTimeout = time.Duration(timeOut) * time.Second
}
