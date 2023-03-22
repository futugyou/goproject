package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/ServiceWeaver/weaver"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

//go:generate weaver generate ./...

func main() {
	// Get a network listener on address "localhost:12345".
	root := weaver.Init(context.Background())
	opts := weaver.ListenerOptions{LocalAddress: "localhost:12345"}
	lis, err := root.Listener("hello", opts)
	if err != nil {
		root.Logger().Error("Listener localhost:12345", err)
	}

	root.Logger().Info(fmt.Sprintf("hello listener available on %v\n", lis))

	// Get a client to the Reverser component.
	reverser, err := weaver.Get[Reverser](root)
	if err != nil {
		log.Fatal(err)
	}

	// Serve the /hello endpoint.
	http.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		reversed, err := reverser.Reverse(r.Context(), r.URL.Query().Get("name"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		fmt.Fprintf(w, "Hello, %s!\n", reversed)
	})

	// Get a client to the IModel component.
	model, err := weaver.Get[IModel](root)
	if err != nil {
		log.Fatal(err)
	}

	// Serve the /model endpoint.
	http.HandleFunc("/model", func(w http.ResponseWriter, r *http.Request) {
		models, err := model.ListModel(r.Context())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		w.Header().Set("content-type", "text/json")
		msg, _ := json.Marshal(models)
		w.Write(msg)
		trace.SpanFromContext(r.Context()).AddEvent("writing response",
			trace.WithAttributes(
				attribute.String("content", "hello "),
				attribute.String("answer", r.URL.Query().Get("name")),
			))
	})

	// Create an otel handler to enable tracing.
	otelHandler := otelhttp.NewHandler(http.DefaultServeMux, "http")
	http.Serve(lis, otelHandler)
}
