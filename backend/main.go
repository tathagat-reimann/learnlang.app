package main

import (
	"log"
	"net/http"

	"learnlang-backend/router"
	"learnlang-backend/store"
	"learnlang-backend/utils"
)

func main() {
	if err := store.InitFromEnv(); err != nil {
		log.Fatalf("failed to init database: %v", err)
	}
	if err := utils.VerifyUploadDirWritable(); err != nil {
		log.Fatalf("upload dir check failed: %v", err)
	}
	r := router.NewRouter()

	log.Println("Server running on :8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatal(err)
	}
}
