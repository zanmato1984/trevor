package handlers

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

func WorkHandler(w http.ResponseWriter, r *http.Request) {
	tr := otel.Tracer("go-demo-server")
	ctx, span := tr.Start(r.Context(), "work-handler")
	defer span.End()

	// Add useful span attributes
	span.SetAttributes(
		attribute.String("http.method", r.Method),
		attribute.String("http.url", r.URL.String()),
		attribute.String("user_agent", r.UserAgent()),
		attribute.String("client_ip", r.RemoteAddr),
		attribute.String("custom.demo_id", "12345"),
		attribute.String("custom.demo_type", "enriched-demo"),
	)

	// Simulate two overlapping child operations (e.g., DB query and external API call)
	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		_, dbSpan := tr.Start(ctx, "db-query")
		dbSpan.SetAttributes(
			attribute.String("db.system", "postgresql"),
			attribute.String("db.statement", "SELECT * FROM demo WHERE id=12345"),
			attribute.Int("db.rows_returned", 1),
			attribute.String("db.user", "demo_user"),
			attribute.String("db.operation", "SELECT"),
			attribute.String("db.success", "true"),
		)
		dbSpan.AddEvent("DB query started")
		time.Sleep(50 * time.Millisecond)
		// Simulate an error
		err := fmt.Errorf("simulated DB error")
		dbSpan.RecordError(err)
		dbSpan.SetStatus(codes.Error, err.Error())
		dbSpan.AddEvent("DB query failed", trace.WithAttributes(attribute.String("error", err.Error())))
		dbSpan.End()
	}()

	go func() {
		defer wg.Done()
		_, apiSpan := tr.Start(ctx, "external-api-call")
		apiSpan.SetAttributes(
			attribute.String("http.method", "GET"),
			attribute.String("http.url", "https://api.example.com/resource"),
			attribute.Int("http.status_code", 200),
			attribute.String("api.version", "v1"),
			attribute.String("api.caller", "work-handler"),
			attribute.String("api.success", "true"),
		)
		apiSpan.AddEvent("API call started")
		time.Sleep(30 * time.Millisecond)
		apiSpan.AddEvent("API call finished", trace.WithAttributes(attribute.Int("http.status_code", 200)))
		apiSpan.End()
	}()

	wg.Wait()

	// Simulate work
	fmt.Fprintln(w, "Work done with child spans!")
}
