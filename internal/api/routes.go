package api

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/arun-builds/gridfall/internal/api/admin"
	"github.com/arun-builds/gridfall/internal/api/user"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func (a *Api) RegisterRoutes() http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	r.Get("/", a.HelloWorldHandler)

	userHandler := user.NewHandler(a.userRepo)
	user.RegisterRoutes(r, userHandler)

	adminHandler := admin.NewHandler(a.adminRepo)
	admin.RegisterRoutes(r, adminHandler)

	return r
}

func (s *Api) HelloWorldHandler(w http.ResponseWriter, r *http.Request) {
	resp := make(map[string]string)
	resp["message"] = "Hello World"

	jsonResp, err := json.Marshal(resp)
	if err != nil {
		log.Fatalf("error handling JSON marshal. Err: %v", err)
	}

	_, _ = w.Write(jsonResp)
}
