package main

import (
	"github.com/joho/godotenv"
	"github.com/tools"
	"os"
)

var BaseURL string

func init() {
	godotenv.Load()
	BaseURL = os.Getenv("BaseURL")
}

func main() {
	var client tools.Client
	client.BaseURL = BaseURL

	tools.HistoricalPrice(client)
}
