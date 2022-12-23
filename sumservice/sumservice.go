package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/hsmtkk/jubilant-adventure/env"
)

func main() {
	port, err := env.Port()
	if err != nil {
		log.Fatal(err)
	}

	// Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Routes
	e.GET("/", hello)
	e.POST("/sum", sum)

	// Start server
	e.Logger.Fatal(e.Start(fmt.Sprintf(":%d", port)))
}

// Handler
func hello(c echo.Context) error {
	return c.String(http.StatusOK, "Hello Challenger01!")
}

type sumRequestSchema struct {
	Numbers []int `json:"numbers"`
}

type sumResponseSchema struct {
	Sum int `json:"sum"`
}

func sum(ectx echo.Context) error {
	sumReq := new(sumRequestSchema)
	if err := ectx.Bind(sumReq); err != nil {
		return fmt.Errorf("echo.Contex.Bind failed; %w", err)
	}
	result := 0
	for _, n := range sumReq.Numbers {
		result += n
	}
	return ectx.JSON(http.StatusOK, sumResponseSchema{Sum: result})
}
