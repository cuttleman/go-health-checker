package main

import (
	"healthChecker"
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
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

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

	e.Logger.Fatal(e.Start(":4000"))
}
