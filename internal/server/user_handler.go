package server

import (
	"encoding/json"
	"net/http"

	"github.com/dakshcodez/gdg_chat_app_backend_task/internal/middleware"
	"github.com/dakshcodez/gdg_chat_app_backend_task/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserHandler struct {
	DB *gorm.DB
}

func (h *UserHandler) Me(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserIDKey)

	var user models.User
	if err := h.DB.First(&user, "id = ?", userID).Error; err != nil {
		http.Error(w, "user not found", http.StatusNotFound)
		return
	}

	// Return safe fields only
	response := map[string]interface{}{
		"id":       user.ID,
		"username": user.Username,
		"email":    user.Email,
		"created":  user.CreatedAt,
	}

	json.NewEncoder(w).Encode(response)
}

func (h *UserHandler) UpdateMe(w http.ResponseWriter, r *http.Request) {
	// Get authenticated user ID
	userID := r.Context().Value(middleware.UserIDKey).(uuid.UUID)

	// Parse request body
	var body struct {
		Username *string `json:"username"`
		Email    *string `json:"email"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	// Build update map (partial updates allowed)
	updates := map[string]interface{}{}

	if body.Username != nil {
		updates["username"] = *body.Username
	}

	if body.Email != nil {
		updates["email"] = *body.Email
	}

	if len(updates) == 0 {
		http.Error(w, "no fields to update", http.StatusBadRequest)
		return
	}

	// Perform update (authorized by userID)
	if err := h.DB.
		Model(&models.User{}).
		Where("id = ?", userID).
		Updates(updates).Error; err != nil {

		http.Error(w, "update failed", http.StatusBadRequest)
		return
	}

	// Fetch updated user
	var user models.User
	h.DB.First(&user, "id = ?", userID)

	// Return safe response
	response := map[string]interface{}{
		"id":       user.ID,
		"username": user.Username,
		"email":    user.Email,
		"updated":  user.UpdatedAt,
	}

	json.NewEncoder(w).Encode(response)
}
