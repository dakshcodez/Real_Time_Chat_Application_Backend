package auth

import (
	"encoding/json"
	"net/http"

	"gorm.io/gorm"
)

type Handler struct {
	DB *gorm.DB
	Secret string
}

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Username string
		Email    string
		Password string 
	}
	json.NewDecoder(r.Body).Decode(&body)

	err := Register(h.DB, body.Username, body.Email, body.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Email    string
		Password string
	}
	json.NewDecoder(r.Body).Decode(&body)

	user, err := Login(h.DB, body.Email, body.Password)
	if err != nil {
		http.Error(w, "Invalid Credentials", http.StatusUnauthorized)
		return
	}

	token, _ := GenerateJWT(user.ID, h.Secret)

	json.NewEncoder(w).Encode(map[string]string{
		"token" : token,
	})
}