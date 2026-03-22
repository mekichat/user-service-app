package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"user-service-app/internals/service"
)

type UserHandler struct {
	userService *service.UserService
}

func NewUserHandler(s *service.UserService) *UserHandler {
	return &UserHandler{userService: s}
}

func (h *UserHandler) RegisterRoutes() {
	http.HandleFunc("/users", h.handleUsers)
	http.HandleFunc("/users/", h.handleUsersByID)
}

// /users/{id}

func (h *UserHandler) handleUsersByID(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-type", "application/json")

	idStr := strings.TrimPrefix(r.URL.Path, "/users/")
	id, err := strconv.Atoi(idStr)

	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodDelete:
		ok := h.userService.DeleteUser(id)

		if !ok {
			http.Error(w, "Cant find the user", http.StatusNotFound)
			return
		}

		w.WriteHeader(http.StatusNoContent) // 204

	case http.MethodPut:
		var input struct {
			Name string `json:"name"`
			Age  int    `json:"age"`
		}

		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			http.Error(w, "Invalid json", http.StatusBadRequest)
			return
		}

		user, ok := h.userService.UpdateUser(id, input.Name, input.Age)

		if !ok {
			http.Error(w, "Can not find the user", http.StatusNotFound)
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(user)

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)

	}

}

func (h *UserHandler) handleUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case http.MethodGet:
		users := h.userService.ListUsers()
		json.NewEncoder(w).Encode(users)

	case http.MethodPost:

		var input struct {
			Name string `json:"name"`
			Age  int    `json:"age"`
		}

		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		user := h.userService.CreateUser(input.Name, input.Age)

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(user)

	default:
		http.Error(w, "Invalid Method", http.StatusMethodNotAllowed)
	}

}
