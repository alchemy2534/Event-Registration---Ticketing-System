package handlers

import (
	"context"
	"encoding/json"
	"event-registration-system/internal/models"
	"event-registration-system/internal/repository"
	"net/http"
)

type RegistrationHandler struct {
	repo *repository.RegistrationRepository
}

func NewRegistrationHandler(repo *repository.RegistrationRepository) *RegistrationHandler {
	return &RegistrationHandler{repo: repo}
}

func (h *RegistrationHandler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var reg models.Registration
	if err := json.NewDecoder(r.Body).Decode(&reg); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	err := h.repo.RegisterForEvent(context.Background(), &reg)

	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		w.WriteHeader(http.StatusConflict) // Conflict or bad request if already registered/full
		json.NewEncoder(w).Encode(models.RegistrationResponse{
			Message: "Registration failed: " + err.Error(),
			Status:  false,
		})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(models.RegistrationResponse{
		Message: "Registered successfully",
		Status:  true,
	})
}
