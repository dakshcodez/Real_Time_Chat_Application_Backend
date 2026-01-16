package server

import (
	"encoding/json"
	"net/http"

	"github.com/dakshcodez/gdg_chat_app_backend_task/internal/middleware"
	"github.com/dakshcodez/gdg_chat_app_backend_task/internal/models"
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
