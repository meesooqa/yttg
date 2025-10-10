package main

import (
	"context"
	"log"

	"github.com/joho/godotenv"

	"github.com/meesooqa/yttg/app/job"
	"github.com/meesooqa/yttg/app/web"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	jq := job.NewJobQueue()

	// Start workers
	numWorkers := 1
	for i := 1; i <= numWorkers; i++ {
		go job.Worker(i, jq)
	}

	server := web.Server{
		JobQueue: jq,
	}
	server.Run(context.Background(), 8080)
}
