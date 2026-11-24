package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync/atomic"
)

var counter atomic.Uint64

func podName() string {
	if n := os.Getenv("POD_NAME"); n != "" {
		return n
	}
	if n, err := os.Hostname(); err == nil {
		return n
	}
	return "unknown"
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func main() {
	pod := podName()

	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
	})

	mux.HandleFunc("GET /id", func(w http.ResponseWriter, r *http.Request) {
		id := counter.Add(1)
		writeJSON(w, http.StatusOK, map[string]string{
			"id":           strconv.FormatUint(id, 10),
			"generated_by": pod,
		})
	})

	addr := ":8080"
	log.Printf("notiflex-api starting on %s (pod=%s)", addr, pod)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatal(err)
	}
}
