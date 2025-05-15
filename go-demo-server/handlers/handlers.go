package handlers

import (
	"fmt"
	"net/http"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

func WorkHandler(w http.ResponseWriter, r *http.Request) {
	tr := otel.Tracer("go-demo-server")
	_, span := tr.Start(r.Context(), "work-handler")
	defer span.End()

	// Add useful span attributes
	span.SetAttributes(
		attribute.String("http.method", r.Method),
		attribute.String("http.url", r.URL.String()),
		attribute.String("user_agent", r.UserAgent()),
		attribute.String("client_ip", r.RemoteAddr),
	)

	// Simulate work
	fmt.Fprintln(w, "Work done!")
}
