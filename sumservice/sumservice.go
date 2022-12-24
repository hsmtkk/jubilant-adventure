package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"google.golang.org/api/idtoken"

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
	e.POST("/sumcurrency", sumCurrency)

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

type sumCurrencyRequestSchema struct {
	Amounts []string `json:"amounts"`
}

type sumCurrencyResponseSchema struct {
	Sum int `json:"sum"`
}

func sumCurrency(ectx echo.Context) error {
	req := new(sumCurrencyRequestSchema)
	if err := ectx.Bind(req); err != nil {
		return fmt.Errorf("echo.Context.Bind failed; %w", err)
	}
	sum, err := invokeCurrencyService(ectx.Request().Context(), req.Amounts)
	if err != nil {
		return err
	}
	return ectx.JSON(http.StatusOK, sumCurrencyResponseSchema{Sum: sum})
}

type convertRequestSchema struct {
	Value string `json:"value"`
}

type convertResponseSchema struct {
	Answer int `json:"answer"`
}

func invokeCurrencyService(ctx context.Context, amounts []string) (int, error) {
	audience, err := env.RequiredVar("CURRENCY_SERVICE")
	if err != nil {
		return 0, err
	}
	url := audience + "/convert"
	clt, err := idtoken.NewClient(ctx, audience)
	if err != nil {
		return 0, fmt.Errorf("idtoken.NewClient failed; %w", err)
	}
	sum := 0
	for _, amount := range amounts {
		reqBytes, err := json.Marshal(convertRequestSchema{
			Value: amount,
		})
		if err != nil {
			return 0, fmt.Errorf("json.Marshal failed; %w", err)
		}
		req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(reqBytes))
		if err != nil {
			return 0, fmt.Errorf("http.NewRequest failed; %w", err)
		}
		resp, err := clt.Do(req)
		if err != nil {
			return 0, fmt.Errorf("http.Client.Do failed; %w", err)
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			return 0, fmt.Errorf("non 200 HTTP response; %d; %s", resp.StatusCode, resp.Status)
		}
		respBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			return 0, fmt.Errorf("io.ReadAll failed; %w", err)
		}
		var respSchema convertResponseSchema
		if err := json.Unmarshal(respBytes, &respSchema); err != nil {
			return 0, fmt.Errorf("json.Unmarshal failed; %w", err)
		}
		sum += respSchema.Answer
	}
	return sum, nil
}
