package main

import (
	"log"
	"net/http"
	"slices"
	"time"
)

type loggingResponseWriter struct {
	http.ResponseWriter
	status int
}

func loggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		lrw := &loggingResponseWriter{ResponseWriter: w, status: http.StatusOK}

		next.ServeHTTP(lrw, r)

		log.Printf("%s %s %d %v",
			r.Method, r.URL.Path, lrw.status, time.Since(start))
	})
}

func recoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("panic: %v", err)
				writeError(w, http.StatusInternalServerError, "internal server error")
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func chain(h http.Handler, mw ...func(http.Handler) http.Handler) http.Handler {
	// for i := len(mw) - 1; i >= 0; i-- {
	// 	h = mw[i](h)
	// }
	for _, m := range slices.Backward(mw) {
		h = m(h)
	}
	return h
}
