package main

import (
	"log"

	"rip-go-app/internal/api"
)

func main() {
	log.Println("Application start up")
	api.StartServer()
	log.Println("Application terminated")
}