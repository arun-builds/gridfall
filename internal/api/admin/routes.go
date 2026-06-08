package admin

import "github.com/go-chi/chi/v5"

func RegisterRoutes(r chi.Router, h *Handler) {
	r.Route("/admin", func(r chi.Router) {
		r.Get("/users", h.ListUsers)
		r.Get("/users/{id}", h.GetUser)
		r.Put("/users/{id}/role", h.UpdateUserRole)
		r.Delete("/users/{id}", h.DeleteUser)
	})
}
