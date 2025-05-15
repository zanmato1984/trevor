package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"go-demo-server/handlers"
	"go-demo-server/tracing"
)

func main() {
	ctx := context.Background()

	tracerProvider, err := tracing.InitTracer(ctx)
	if err != nil {
		log.Fatalf("failed to create exporter: %v", err)
	}
	defer func() {
		if tracerProvider != nil {
			_ = tracerProvider.Shutdown(ctx)
		}
	}()

	http.HandleFunc("/work", handlers.WorkHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Listening on :%s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
