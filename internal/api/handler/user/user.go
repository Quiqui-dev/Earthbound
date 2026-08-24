package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/Quiqui-dev/Earthbound/internal/model"
	service "github.com/Quiqui-dev/Earthbound/internal/service/user"
	"github.com/Quiqui-dev/Earthbound/util"
)

type UserHandler struct {
	service *service.UserService
}

func NewUserHandler(service *service.UserService) *UserHandler {
	return &UserHandler{
		service: service,
	}
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var req model.RequestCreateUser
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("Create User - JSON Decode error: %v", err)
		util.RespondWithErrorFromCode(w, http.StatusBadRequest)
		return
	}

	log.Printf("Create User - Request received: username=%s email=%s", req.Username, req.Email)

	user, err := h.service.CreateUser(r.Context(), req)

	if err != nil {
		log.Printf("Create User - Service Error: %v", err)
		util.RespondWithErrorCustom(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Printf("Create User - Success: user created with id=%s, username=%s", user.ID, user.Username)

	http.SetCookie(w, &http.Cookie{
		Name:     "jwt",
		Value:    user.AccessToken,
		Path:     "/",
		MaxAge:   60 * 60 * 24,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	})

	util.WriteJSON(w, http.StatusCreated, user)
}

func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req model.RequestLoginUser
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		util.RespondWithErrorFromCode(w, http.StatusBadRequest)
		return
	}

	user, err := h.service.Login(r.Context(), req)
	if err != nil {
		util.RespondWithErrorCustom(w, err.Error(), http.StatusUnauthorized)
		return
	}

	// Set JWT cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "jwt",
		Value:    user.AccessToken,
		Path:     "/",
		MaxAge:   60 * 60 * 24,
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
	})

	util.WriteJSON(w, http.StatusOK, user)
}

func (h *UserHandler) Logout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:     "jwt",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
	})

	util.WriteJSON(w, http.StatusOK, map[string]string{"message": "logout successful"})
}
