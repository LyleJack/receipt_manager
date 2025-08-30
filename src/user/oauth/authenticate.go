package oauth

import (
	_ "context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/idtoken"
)

var (
	oauthConfig = &oauth2.Config{
		ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		RedirectURL:  os.Getenv("GOOGLE_REDIRECT_URL"),
		Scopes:       []string{"openid", "profile", "email"},
		Endpoint:     google.Endpoint,
	}
)

func generateOAuthState() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

func AuthenticationURL(w http.ResponseWriter, r *http.Request) (string, error) {
	state, err := generateOAuthState()
	if err != nil {
		return "", err
	}

	// Store state in secure cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "oauthstate",
		Value:    state,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
		MaxAge:   600, // seconds
	})

	return oauthConfig.AuthCodeURL(state, oauth2.AccessTypeOffline), nil
}

// HandleOAuthCallback processes the OAuth2 callback from Google
func HandleOAuthCallback(w http.ResponseWriter, r *http.Request) {
	stateCookie, err := r.Cookie("oauthstate")
	if err != nil || r.FormValue("state") != stateCookie.Value {
		http.Error(w, "Invalid state parameter", http.StatusBadRequest)
		return
	}

	// Exchange code for token
	code := r.FormValue("code")
	token, err := oauthConfig.Exchange(r.Context(), code)
	if err != nil {
		log.Printf("OAuth exchange failed: %v", err)
		http.Error(w, "Login failed", http.StatusInternalServerError)
		return
	}

	// Validate ID token
	rawIDToken, ok := token.Extra("id_token").(string)
	if !ok {
		http.Error(w, "No ID token in response", http.StatusInternalServerError)
		return
	}

	payload, err := idtoken.Validate(r.Context(), rawIDToken, oauthConfig.ClientID)
	if err != nil {
		log.Printf("ID token validation failed: %v", err)
		http.Error(w, "Login failed", http.StatusUnauthorized)
		return
	}

	// Extract user info
	userInfo := map[string]interface{}{
		"email":     payload.Claims["email"],
		"name":      payload.Claims["name"],
		"picture":   payload.Claims["picture"],
		"google_id": payload.Claims["sub"],
		"verified":  payload.Claims["email_verified"],
	}

	// TODO: Find or create user in your DB here
	// user := db.FindOrCreateUser(userInfo)

	// TODO: Start a session or set a JWT
	// startSession(w, user)

	// For now, just return JSON response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(userInfo)
}
