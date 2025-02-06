package main

import (
	"context"

	"task-queue-001/app/job"
	"task-queue-001/app/web"
)

func main() {
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
