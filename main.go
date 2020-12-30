package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

func main() {

	s := NewServer("Olis Server 1.1", getCounterfile())

	s.logger.Printf("Server is starting on %s...", getServerPort())
	s.logger.Printf("Counter file: %s...", getCounterfile())

	s.setupRoutes()
	
	http.ListenAndServe(getServerPort(), s.router)

}

func (s *server) setupRoutes() {
	s.router.HandleFunc("/", s.logRequest(s.handleIndex()))
	s.router.HandleFunc("/showlog", s.logRequest(s.handleShowLog()))
	s.router.NotFoundHandler = s.handleNotFound()
}

func (s *server) handleIndex() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		hello := fmt.Sprintf("Hello World! from %s", s.serverName)
		w.Write([]byte(hello))
	}
}

func (s *server) handleShowLog() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(s.logInfo)
	}
}

func getCounterfile() string {
	if filename, ok := os.LookupEnv("COUNTERFILE"); ok {
		return filename
	}
	return "counter.json"
}

func getServerPort() string {
	if port, ok := os.LookupEnv("PORT"); ok {
		return ":" + port
	}
	return ":3000"
}