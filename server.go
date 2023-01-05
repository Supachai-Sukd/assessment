package main

import (
	"context"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "github.com/lib/pq"
	"github.com/supachai-sukd/assessment/app/customer_expense"
	"github.com/supachai-sukd/assessment/pkg/config"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func main() {
	config.InitDB()

	e := echo.New()

	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{echo.GET, echo.HEAD, echo.PUT, echo.PATCH, echo.POST, echo.DELETE},
	}))

	e.GET("/", func(c echo.Context) error {
		return c.JSON(http.StatusOK, "OK")
	})

	e.POST("/expenses", customer_expense.AddExpenses)
	e.GET("/expenses/:id", customer_expense.GetExpensesById)
	e.PUT("/expenses/:id", customer_expense.UpdateExpenses)
	e.GET("/expenses", customer_expense.GetAllExpenses)

	defaultPORT := os.Getenv("PORT")
	if defaultPORT == "" {
		defaultPORT = "2565"
	}

	go func() {
		if err := e.Start(":" + defaultPORT); err != nil && err != http.ErrServerClosed { // Start server
			e.Logger.Fatal("shutting down the server")
		}
	}()
	fmt.Printf("start at port: %v\n", defaultPORT) // Port 2565

	// Create channel shutdown
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt)

	<-shutdown

	// Shutdown in 10 second
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}

}
