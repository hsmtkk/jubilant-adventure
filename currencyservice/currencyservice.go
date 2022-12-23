package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

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
	e.GET("/", healthz)
	e.POST("/convert", convert)

	// Start server
	e.Logger.Fatal(e.Start(fmt.Sprintf(":%d", port)))
}

// Handler
func healthz(c echo.Context) error {
	return c.String(http.StatusOK, "OK")
}

type convertRequestSchema struct {
	Value string `json:"value"`
}

type convertResponseSchema struct {
	Answer int `json:"answer"`
}

func convert(ectx echo.Context) error {
	convReq := new(convertRequestSchema)
	if err := ectx.Bind(convReq); err != nil {
		return fmt.Errorf("echo.Contex.Bind failed; %w", err)
	}
	unit := convReq.Value[:3]
	amount, err := strconv.Atoi(convReq.Value[3:])
	if err != nil {
		return fmt.Errorf("failed to parse as int; %s; %w", convReq.Value[3:], err)
	}
	yen, err := toYen(unit)
	if err != nil {
		return err
	}
	return ectx.JSON(http.StatusOK, convertResponseSchema{Answer: amount * yen})
}

func toYen(unit string) (int, error) {
	switch unit {
	case "JPY":
		return 1, nil
	case "USD":
		return 132, nil
	case "EUR":
		return 140, nil
	case "GBP":
		return 158, nil
	case "CHF":
		return 142, nil
	case "AUD":
		return 87, nil
	case "NZD":
		return 82, nil
	default:
		return 0, fmt.Errorf("unknown unit; %s", unit)
	}
}
