package main

import (
	"fmt"
	"healthChecker"
	"time"
)

func main() {
	start := time.Now()
	greatNode, err := healthChecker.Execute(97)
	since := time.Since(start)

	if err != nil {
		fmt.Println("HealthChecker Error :", err)

		return
	}
	fmt.Println("greatNode :", greatNode)
	fmt.Println(since)
}
