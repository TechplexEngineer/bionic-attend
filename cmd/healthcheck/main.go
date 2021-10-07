package main

import (
	"fmt"
	"net/http"
	"os"
)

// Simple healthcheck tool for from-scratch containers to support AWS ECS healthchecks
// see: https://github.com/Soluto/golang-docker-healthcheck-example

func main() {
	_, err := http.Get(fmt.Sprintf("http://127.0.0.1:%s/", os.Getenv("PORT")))
	if err != nil {
		os.Exit(1)
	}
}