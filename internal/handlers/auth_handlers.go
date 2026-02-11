package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"final-by-me/internal/auth"
	"final-by-me/internal/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type AuthHandler struct {
	users     *mongo.Collection
	jwtSecret []byte
}

func NewAuthHandler(db *mongo.Database, jwtSecret []byte) *AuthHandler {
	return &AuthHandler{
		users:     db.Collection("user"),
		jwtSecret: jwtSecret,
	}
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email    string `json:"email"`
		Name     string `json:"name"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, 400, map[string]string{"error": "invalid JSON"})
		return
	}

	req.Email = strings.ToLower(strings.TrimSpace(req.Email))
	req.Name = strings.TrimSpace(req.Name)

	if req.Email == "" || req.Password == "" {
		writeJSON(w, 400, map[string]string{"error": "email and password required"})
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	// Check existing
	count, err := h.users.CountDocuments(ctx, bson.M{"email": req.Email})
	if err != nil {
		writeJSON(w, 500, map[string]string{"error": "db error"})
		return
	}
	if count > 0 {
		writeJSON(w, 409, map[string]string{"error": "email already exists"})
		return
	}

	hash, err := auth.HashPassword(req.Password)
	if err != nil {
		writeJSON(w, 500, map[string]string{"error": "hash error"})
		return
	}

	u := models.User{
		Email:        req.Email,
		Name:         req.Name,
		PasswordHash: hash,
		Role:         "fan",
		CreatedAt:    time.Now(),
	}

	res, err := h.users.InsertOne(ctx, u)
	if err != nil {
		writeJSON(w, 500, map[string]string{"error": "insert error"})
		return
	}

	writeJSON(w, 201, map[string]any{
		"id":    res.InsertedID,
		"email": u.Email,
		"name":  u.Name,
		"role":  u.Role,
	})
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, 400, map[string]string{"error": "invalid JSON"})
		return
	}

	req.Email = strings.ToLower(strings.TrimSpace(req.Email))

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	var u models.User
	if err := h.users.FindOne(ctx, bson.M{"email": req.Email}).Decode(&u); err != nil {
		writeJSON(w, 401, map[string]string{"error": "invalid credentials"})
		return
	}

	if !auth.CheckPassword(u.PasswordHash, req.Password) {
		writeJSON(w, 401, map[string]string{"error": "invalid credentials"})
		return
	}

	token, err := auth.Sign(h.jwtSecret, u.ID.Hex(), u.Email, u.Role)
	if err != nil {
		writeJSON(w, 500, map[string]string{"error": "token error"})
		return
	}

	writeJSON(w, 200, map[string]any{
		"token": token,
		"role":  u.Role,
	})
}
