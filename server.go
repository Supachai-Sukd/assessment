package main

import (
	"context"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	"github.com/supachai-sukd/assessment/pkg/database"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func main() {
	api := echo.New()
	api.Logger.SetLevel(log.INFO)
	api.Use(middleware.Recover())
	api.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{echo.GET, echo.HEAD, echo.PUT, echo.PATCH, echo.POST, echo.DELETE},
	}))
	db := database.GetInstance()
	if db == nil {
		log.Fatalf("Could not connect to database")
	}
	defer db.Close()

	sqlFile, sqlErr := ioutil.ReadFile("platform/migrations/create_init_tables.up.sql")
	if sqlErr != nil {
		log.Fatal(sqlErr)
	}

	_, err := db.Exec(string(sqlFile))
	if err != nil {
		log.Fatal(err)
	}

	api.GET("/", func(c echo.Context) error {
		return c.JSON(http.StatusOK, "OK")
	})

	// os.Getenv("PORT") Use after refactor.
	go func() {
		if err := api.Start(":2565"); err != nil && err != http.ErrServerClosed { // Start server
			api.Logger.Fatal("shutting down the server")
		}
	}()
	fmt.Printf("start at port: %v\n", 2565) // Port 2565

	// Create channel shutdown
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt)

	<-shutdown

	// Shutdown in 10 second
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := api.Shutdown(ctx); err != nil {
		api.Logger.Fatal(err)
	}

}
