package handler

import (
	"encoding/json"
	"net/http"
	"stock-api/internal/service"
	"strconv"
)

type UserHandler struct {
    userService *service.UserService
}

func NewUserHandler(us *service.UserService) *UserHandler {
    return &UserHandler{userService: us}
}


type Credentials struct {
    Username string `json:"username"`
    Password string `json:"password"`
}

func (u *UserHandler)Register(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }
	var creds Credentials
    json.NewDecoder(r.Body).Decode(&creds)
	err := u.userService.RegisterUser(creds.Username, creds.Password)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
    w.WriteHeader(http.StatusCreated)
}

func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	var creds Credentials
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	token, userId, err := h.userService.LoginUser(creds.Username, creds.Password)
	if err != nil {
		if err.Error() == "user not found" || err.Error() == "invalid password" {
			http.Error(w, err.Error(), http.StatusUnauthorized)
		} else {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"token": token, "userId": strconv.Itoa(userId),})
}
