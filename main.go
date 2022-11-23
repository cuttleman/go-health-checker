package main

import (
	"fmt"
	healthChecker "healthChecker_"
	"os"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Result struct {
	Status       string `json:"status"`
	GreatNode    string `json:"greatNode"`
	ErrorMessage string `json:"errorMessage"`
}

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		port = "4000"
	}

	e := echo.New()

	e.Use(middleware.Logger())

	e.GET("/", func(c echo.Context) error {
		return c.String(200, "Hello Health Checker")
	})

	e.GET("/chain/:id", func(c echo.Context) error {
		chainId := c.Param("id")
		chainIdToInt, _ := strconv.Atoi(chainId)

		// Health checker
		start := time.Now()
		greatNode, err := healthChecker.Execute(uint64(chainIdToInt))
		responseTime := time.Since(start)
		fmt.Println("HealthChecker Response Time :", responseTime)

		result := new(Result)
		statusCode := 200

		if err != nil {
			statusCode = 400
			result.Status = "fail"
			result.ErrorMessage = err.Error()
		} else {
			result.Status = "ok"
			result.GreatNode = greatNode
		}

		return c.JSON(statusCode, result)
	})

	e.Logger.Fatal(e.Start(":" + port))
}
