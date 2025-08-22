package main

import (
	"log"
	"net/http"

	"learnlang-backend/router"
)

func main() {
	r := router.NewRouter()

	log.Println("Server running on :8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatal(err)
	}
}
