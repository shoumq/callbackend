package handlers

import (
	"encoding/json"
	"net/http"
	"vasek/internal/services"
	"golang.org/x/text/encoding/charmap"
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

    // Получаем данные из формы
    name := r.FormValue("name")
    text := r.FormValue("text")
    phone := r.FormValue("phone")

    if name == "" || text == "" {
        http.Error(w, "Missing required fields", http.StatusBadRequest)
        return
    }

    // Попытка перекодировать из Windows-1251 в UTF-8
    if utf8Name, err := charmap.Windows1251.NewDecoder().String(name); err == nil {
        name = utf8Name
    }
    if utf8Text, err := charmap.Windows1251.NewDecoder().String(text); err == nil {
        text = utf8Text
    }

    // Создание запроса через сервис
    id, err := h.requestService.CreateRequest(name, text, phone)
    if err != nil {
        http.Error(w, "Failed to create request: "+err.Error(), http.StatusInternalServerError)
        return
    }

    // Формируем ответ
    response := map[string]interface{}{
        "id":      id,
        "message": "Request created successfully",
    }

    w.Header().Set("Content-Type", "application/json; charset=utf-8")
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(response)
}

func fixEncoding(input string) (string, error) {
    decoder := charmap.Windows1251.NewDecoder()
    output, err := decoder.String(input)
    if err != nil {
        return "", err
    }
    return output, nil
}

func (h *RequestHandler) GetRequestsHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodGet {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

    // Получаем заявки из сервиса
    requests, err := h.requestService.GetRequest()
    if err != nil {
        http.Error(w, "Failed to get requests: "+err.Error(), http.StatusInternalServerError)
        return
    }

    // Перекодируем поля из Windows-1251 в UTF-8, если необходимо
    for i := range requests {
        if utf8Name, err := charmap.Windows1251.NewDecoder().String(requests[i].Name); err == nil {
            requests[i].Name = utf8Name
        }
        if utf8Text, err := charmap.Windows1251.NewDecoder().String(requests[i].Text); err == nil {
            requests[i].Text = utf8Text
        }
    }

    // Устанавливаем корректный заголовок JSON с UTF-8
    w.Header().Set("Content-Type", "application/json; charset=utf-8")
    w.WriteHeader(http.StatusOK)

    // Отправляем JSON клиенту
    if err := json.NewEncoder(w).Encode(requests); err != nil {
        http.Error(w, "Failed to encode JSON: "+err.Error(), http.StatusInternalServerError)
        return
    }
}
