package server

import (
	"log"
	"net/http"
	"vasek/internal/handlers"
)

type Server struct {
	clientHandler *handlers.RequestHandler
}

func NewServer(clientHandler *handlers.RequestHandler) *Server {
	return &Server{clientHandler: clientHandler}
}

func enableCORS(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	}
}

func (s *Server) Start() {
	http.HandleFunc("/get", enableCORS(s.clientHandler.GetRequestsHandler))
	http.HandleFunc("/create", enableCORS(s.clientHandler.CreateRequestHandler))

	log.Println("listen :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
