package api

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/arun-builds/gridfall/internal/api/admin"
	"github.com/arun-builds/gridfall/internal/api/auth"
	"github.com/arun-builds/gridfall/internal/api/user"
	"github.com/arun-builds/gridfall/internal/database/store"
	"github.com/jackc/pgx/v5/pgxpool"
)

type dbPool = *pgxpool.Pool

type Api struct {
	port        int
	dbPool
	queries     *store.Queries
	jwtSecret   string
	userRepo    user.UserRepository
	adminRepo   admin.AdminRepository
	authHandler *auth.Handler
}

func NewApi(port int, pool *pgxpool.Pool) *http.Server {
	queries := store.New(pool)

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "gridfall-dev-secret-change-in-production"
	}

	jwtExpiryHours := 24
	if v := os.Getenv("JWT_EXPIRY_HOURS"); v != "" {
		if parsed, err := strconv.Atoi(v); err == nil {
			jwtExpiryHours = parsed
		}
	}

	a := &Api{
		port:        port,
		dbPool:      pool,
		queries:     queries,
		jwtSecret:   jwtSecret,
		userRepo:    user.NewUserRepository(queries),
		adminRepo:   admin.NewAdminRepository(queries),
		authHandler: auth.NewHandler(queries, jwtSecret, jwtExpiryHours),
	}

	return &http.Server{
		Addr:         fmt.Sprintf(":%d", a.port),
		Handler:      a.RegisterRoutes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}
}

func (a *Api) health(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ok"))
}
