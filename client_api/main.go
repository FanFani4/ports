package main

import (
	"context"
	"fmt"
	"os"

	"github.com/FanFani4/ports/client_api/reader"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

func main() {
	log := logrus.New()
	ctx := context.TODO()

	err := godotenv.Load()
	if err != nil {
		log.Info("no log file")
	}

	jsonPath := os.Getenv("PORTS_JSON_PATH")

	reader, err := reader.NewJSONReader(ctx, log, jsonPath)
	if err != nil {
		log.Fatal(err)
	}

	portsResponse := reader.GetPorts()

	for port := range portsResponse {
		fmt.Println(port)
	}
}
