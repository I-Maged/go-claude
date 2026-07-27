package main

import "net/http"

func buildRouter() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", healthHandler)
	mux.HandleFunc("GET /ping", pingHandler)
	mux.HandleFunc("POST /echo", echoHandler)
	mux.HandleFunc("GET /panic", panicHandler) // demo recovery middleware
	mux.HandleFunc("/", notFoundHandler)

	return chain(mux, loggerMiddleware, recoveryMiddleware)
}
