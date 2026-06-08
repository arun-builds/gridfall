package api

import (
	"fmt"
	"net/http"
	"time"

	"github.com/arun-builds/gridfall/internal/api/admin"
	"github.com/arun-builds/gridfall/internal/api/user"
	"github.com/arun-builds/gridfall/internal/database/store"
	"github.com/jackc/pgx/v5/pgxpool"
)

type dbPool = *pgxpool.Pool

type Api struct {
	port int
	dbPool
	queries   *store.Queries
	userRepo  user.UserRepository
	adminRepo admin.AdminRepository
}

func NewApi(port int, pool *pgxpool.Pool) *http.Server {
	queries := store.New(pool)

	a := &Api{
		port:      port,
		dbPool:    pool,
		queries:   queries,
		userRepo:  user.NewUserRepository(queries),
		adminRepo: admin.NewAdminRepository(queries),
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
