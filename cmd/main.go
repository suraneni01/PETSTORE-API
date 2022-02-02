package main

import (
	"log"
	"os"
	"pet-store-server/internal/petstore"
	"pet-store-server/internal/server"
	"time"
)

func main() {
	var (
		petstoreURL     = getEnvOrDefault("PETSTORE_URL", "http://petstore-demo-endpoint.execute-api.com/petstore")
		petstoreTimeout = shouldParseDuration(getEnvOrDefault("PETSTORE_TIMEOUT", "1s"))
		serviceAddr     = getEnvOrDefault("SERVICE_ADDR", ":8080")
	)

	petstore := petstore.NewClient(petstoreURL, petstoreTimeout)
	service := server.NewService(petstore)
	err := service.Start(serviceAddr)
	if err != nil {
		log.Fatalf("server failed: %v\n", err)
	}
}

func getEnvOrDefault(envName string, defaultValue string) string {
	value := os.Getenv(envName)
	if value != "" {
		return value
	}
	return defaultValue
}

func shouldParseDuration(v string) time.Duration {
	d, err := time.ParseDuration(v)
	if err != nil {
		panic(err)
	}
	return d
}
