package auth

import "github.com/go-chi/chi/v5"

func RegisterPublicRoutes(r chi.Router, h *Handler) {
	r.Route("/auth", func(r chi.Router) {
		r.Post("/register", h.Register)
		r.Post("/login", h.Login)
		r.Post("/guest", h.Guest)
		r.Post("/logout", h.Logout)
	})
}

func RegisterProtectedRoutes(r chi.Router, h *Handler) {
	r.Get("/auth/me", h.Me)
	r.Post("/auth/upgrade", h.UpgradeGuest)
}
