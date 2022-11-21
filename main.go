package main

import (
	"healthChecker"
	"log"
	"os"
	"strconv"

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
		log.Fatal("$PORT must be set")
	}

	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/", func(c echo.Context) error {
		return c.String(200, "Hello Health Checker")
	})

	e.GET("/chain/:id", func(c echo.Context) error {
		chainId := c.Param("id")
		chainIdToInt, _ := strconv.Atoi(chainId)

		// Health checker
		greatNode, err := healthChecker.Execute(uint64(chainIdToInt))

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
