package admin

import (
	"context"

	"github.com/arun-builds/gridfall/internal/database/store"
	"github.com/jackc/pgx/v5/pgtype"
)

type AdminRepository interface {
	ListUsers(ctx context.Context) ([]store.User, error)
	GetUser(ctx context.Context, id pgtype.UUID) (store.User, error)
	UpdateUserRole(ctx context.Context, params store.UpdateUserRoleParams) error
	DeleteUser(ctx context.Context, id pgtype.UUID) error
}

type adminRepository struct {
	queries *store.Queries
}

func NewAdminRepository(queries *store.Queries) AdminRepository {
	return &adminRepository{queries: queries}
}

func (r *adminRepository) ListUsers(ctx context.Context) ([]store.User, error) {
	return r.queries.ListUsers(ctx)
}

func (r *adminRepository) GetUser(ctx context.Context, id pgtype.UUID) (store.User, error) {
	return r.queries.GetUser(ctx, id)
}

func (r *adminRepository) UpdateUserRole(ctx context.Context, params store.UpdateUserRoleParams) error {
	return r.queries.UpdateUserRole(ctx, params)
}

func (r *adminRepository) DeleteUser(ctx context.Context, id pgtype.UUID) error {
	return r.queries.DeleteUser(ctx, id)
}
