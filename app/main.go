package main

import (
	"context"

	"github.com/meesooqa/yttg/app/job"
	"github.com/meesooqa/yttg/app/web"
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
