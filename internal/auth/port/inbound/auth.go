package auth

import "net/http"

type AuthHandler interface {
	RegisterPassword(w http.ResponseWriter, r *http.Request)
	Login(w http.ResponseWriter, r *http.Request)
	RefreshToken(w http.ResponseWriter, r *http.Request)
	Logout(w http.ResponseWriter, r *http.Request)
}
