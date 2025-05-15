package handlers

import (
	"fmt"
	"net/http"

	"go.opentelemetry.io/otel"
)

func WorkHandler(w http.ResponseWriter, r *http.Request) {
	tr := otel.Tracer("go-demo-server")
	_, span := tr.Start(r.Context(), "work-handler")
	defer span.End()
	// Simulate work
	fmt.Fprintln(w, "Work done!")
}
