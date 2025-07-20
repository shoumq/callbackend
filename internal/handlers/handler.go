package handlers

import (
	"encoding/json"
	"net/http"
	"vasek/internal/services"
)

type RequestHandler struct {
	requestService *services.RequestService
}

func NewRequestHandler(requestService *services.RequestService) *RequestHandler {
	return &RequestHandler{requestService: requestService}
}

func (h *RequestHandler) CreateRequestHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	name := r.FormValue("name")
	email := r.FormValue("email")
	text := r.FormValue("text")
	phone := r.FormValue("phone")

	if name == "" || email == "" || text == "" {
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}

	id, err := h.requestService.CreateRequest(name, email, text, phone)
	if err != nil {
		http.Error(w, "Failed to create request: "+err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"id":      id,
		"message": "Request created successfully",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

func (h *RequestHandler) GetRequestsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	requests, err := h.requestService.GetRequest()
	if err != nil {
		http.Error(w, "Failed to get requests: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(requests)
}
