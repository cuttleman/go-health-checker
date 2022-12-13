package main

import (
	"fmt"
	"healthchecker-server/internal/chainlist"
	"healthchecker-server/internal/healthchecker"
	"net/http"
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
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{http.MethodGet},
	}))
	e.Static("/assets", "assets")

	e.GET("/", func(c echo.Context) error {
		return c.String(200, "Hello Health Checker\n\n\nRenew Chainlist:\n\t- https://chain-healthchecker.fly.dev/chain/update\n\nGet Healthy RPC:\n\t- https://chain-healthchecker.fly.dev/chain/{chainId}\n\nGet Chainlist JSON:\n\t- https://chain-healthchecker.fly.dev/assets/chainlist.json")
	})

	e.GET("/chain/update", func(c echo.Context) error {
		err := chainlist.Execute()

		if err != nil {
			return c.String(500, "chain list update fail. reason : "+err.Error())
		}

		return c.String(200, "chain list update complete")
	})

	e.GET("/chain/:id", func(c echo.Context) error {
		chainId := c.Param("id")
		chainIdToInt, _ := strconv.Atoi(chainId)

		// Health checker
		start := time.Now()
		greatNode, err := healthchecker.Execute(uint64(chainIdToInt))
		responseTime := time.Since(start)
		fmt.Println("healthChecker Response Time :", responseTime)

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
