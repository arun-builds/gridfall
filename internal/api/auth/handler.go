package auth

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/arun-builds/gridfall/internal/database/store"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"golang.org/x/crypto/bcrypt"
)

type Claims struct {
	UserID      string `json:"user_id"`
	AccountType string `json:"account_type"`
	Role        string `json:"role"`
	jwt.RegisteredClaims
}

const cookieName = "gridfall_token"

type Handler struct {
	queries   *store.Queries
	jwtSecret []byte
	jwtExpiry time.Duration
}

func NewHandler(queries *store.Queries, jwtSecret string, jwtExpiryHours int) *Handler {
	return &Handler{
		queries:   queries,
		jwtSecret: []byte(jwtSecret),
		jwtExpiry: time.Duration(jwtExpiryHours) * time.Hour,
	}
}

func (h *Handler) setTokenCookie(w http.ResponseWriter, token string) {
	http.SetCookie(w, &http.Cookie{
		Name:     cookieName,
		Value:    token,
		Path:     "/",
		MaxAge:   int(h.jwtExpiry.Seconds()),
		HttpOnly: true,
		Secure:   false, // set to true in production with HTTPS
		SameSite: http.SameSiteLaxMode,
	})
}

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.Name == "" || req.Email == "" || req.Password == "" {
		writeError(w, http.StatusBadRequest, "name, email, and password are required")
		return
	}

	if len(req.Password) < 6 {
		writeError(w, http.StatusBadRequest, "password must be at least 6 characters")
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("error hashing password: %v", err)
		writeError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	// TODO: better error handling in case of existing user & error
	user, err := h.queries.CreateRegisteredUser(r.Context(), store.CreateRegisteredUserParams{
		Name:         req.Name,
		Email:        pgtype.Text{String: req.Email, Valid: true},
		PasswordHash: pgtype.Text{String: string(hashedPassword), Valid: true},
	})
	if err != nil {
		log.Printf("error creating user: %v", err)
		writeError(w, http.StatusConflict, "email already registered")
		return
	}

	token, err := h.generateToken(user)
	if err != nil {
		log.Printf("error generating token: %v", err)
		writeError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	h.setTokenCookie(w, token)
	writeJSON(w, http.StatusCreated, map[string]any{
		"user": userResponse(user),
	})
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.Email == "" || req.Password == "" {
		writeError(w, http.StatusBadRequest, "email and password are required")
		return
	}

	user, err := h.queries.GetUserByEmail(r.Context(), pgtype.Text{String: req.Email, Valid: true})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			writeError(w, http.StatusUnauthorized, "invalid email or password")
			return
		}
		log.Printf("error getting user by email: %v", err)
		writeError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	if !user.PasswordHash.Valid {
		writeError(w, http.StatusUnauthorized, "invalid email or password")
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash.String), []byte(req.Password)); err != nil {
		writeError(w, http.StatusUnauthorized, "invalid email or password")
		return
	}

	token, err := h.generateToken(user)
	if err != nil {
		log.Printf("error generating token: %v", err)
		writeError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	h.setTokenCookie(w, token)
	writeJSON(w, http.StatusOK, map[string]any{
		"user": userResponse(user),
	})
}

// Guest handles POST /auth/guest
func (h *Handler) Guest(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Name string `json:"name"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.Name == "" {
		req.Name = "Guest"
	}

	user, err := h.queries.CreateGuestUser(r.Context(), req.Name)
	if err != nil {
		log.Printf("error creating guest user: %v", err)
		writeError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	token, err := h.generateToken(user)
	if err != nil {
		log.Printf("error generating token: %v", err)
		writeError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	h.setTokenCookie(w, token)
	writeJSON(w, http.StatusCreated, map[string]any{
		"user": userResponse(user),
	})
}

func (h *Handler) UpgradeGuest(w http.ResponseWriter, r *http.Request) {
	userID, ok := UserIDFromCtx(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	accountType, _ := AccountTypeFromCtx(r.Context())
	if accountType != string(store.AccountTypeGuest) {
		writeError(w, http.StatusBadRequest, "only guest accounts can be upgraded")
		return
	}

	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.Email == "" || req.Password == "" {
		writeError(w, http.StatusBadRequest, "email and password are required")
		return
	}

	if len(req.Password) < 6 {
		writeError(w, http.StatusBadRequest, "password must be at least 6 characters")
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("error hashing password: %v", err)
		writeError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	var id pgtype.UUID
	if err := id.Scan(userID); err != nil {
		writeError(w, http.StatusBadRequest, "invalid user id in token")
		return
	}

	user, err := h.queries.UpgradeGuestToRegistered(r.Context(), store.UpgradeGuestToRegisteredParams{
		ID:           id,
		Email:        pgtype.Text{String: req.Email, Valid: true},
		PasswordHash: pgtype.Text{String: string(hashedPassword), Valid: true},
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			writeError(w, http.StatusBadRequest, "account is not a guest or does not exist")
			return
		}
		log.Printf("error upgrading guest: %v", err)
		writeError(w, http.StatusConflict, "email already registered")
		return
	}

	token, err := h.generateToken(user)
	if err != nil {
		log.Printf("error generating token: %v", err)
		writeError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	h.setTokenCookie(w, token)
	writeJSON(w, http.StatusOK, map[string]any{
		"user": userResponse(user),
	})
}

func (h *Handler) Me(w http.ResponseWriter, r *http.Request) {
	userIDStr, ok := UserIDFromCtx(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var id pgtype.UUID
	if err := id.Scan(userIDStr); err != nil {
		writeError(w, http.StatusBadRequest, "invalid user id in token")
		return
	}

	user, err := h.queries.GetUser(r.Context(), id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			writeError(w, http.StatusNotFound, "user not found")
			return
		}
		log.Printf("error getting user: %v", err)
		writeError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"user": userResponse(user),
	})
}

func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:     cookieName,
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
	})

	writeJSON(w, http.StatusOK, map[string]string{
		"message": "logged out",
	})
}

func (h *Handler) generateToken(user store.User) (string, error) {
	uuidBytes := user.ID.Bytes
	userID := fmt.Sprintf("%x-%x-%x-%x-%x",
		uuidBytes[0:4], uuidBytes[4:6], uuidBytes[6:8], uuidBytes[8:10], uuidBytes[10:16])

	now := time.Now()
	claims := Claims{
		UserID:      userID,
		AccountType: string(user.Type),
		Role:        string(user.Role),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(h.jwtExpiry)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			Issuer:    "gridfall",
			Subject:   userID,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(h.jwtSecret)
}

func userResponse(u store.User) map[string]any {
	uuidBytes := u.ID.Bytes
	id := fmt.Sprintf("%x-%x-%x-%x-%x",
		uuidBytes[0:4], uuidBytes[4:6], uuidBytes[6:8], uuidBytes[8:10], uuidBytes[10:16])

	resp := map[string]any{
		"id":         id,
		"name":       u.Name,
		"type":       string(u.Type),
		"role":       string(u.Role),
		"created_at": u.CreatedAt.Time,
	}
	if u.Email.Valid {
		resp["email"] = u.Email.String
	}
	return resp
}
