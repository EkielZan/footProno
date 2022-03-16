package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	serverPort := os.Getenv("SERVER_PORT")
	if serverPort == "" {
		serverPort = "8000" //localhost
	}
	_, err = http.Get(fmt.Sprintf("http://127.0.0.1:%s/health", serverPort))
	if err != nil {
		log.Fatalln(err)
		os.Exit(1)
	}
}
