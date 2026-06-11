package user

import (
	"context"

	"github.com/arun-builds/gridfall/internal/database/store"
	"github.com/jackc/pgx/v5/pgtype"
)

type UserRepository interface {
	List(ctx context.Context) ([]store.User, error)
	Get(ctx context.Context, id pgtype.UUID) (store.User, error)
	UpdateName(ctx context.Context, params store.UpdateUserNameParams) error
	UpdateType(ctx context.Context, params store.UpdateUserTypeParams) error
	Delete(ctx context.Context, id pgtype.UUID) error
}

type userRepository struct {
	queries *store.Queries
}

func NewUserRepository(queries *store.Queries) UserRepository {
	return &userRepository{queries: queries}
}

func (r *userRepository) List(ctx context.Context) ([]store.User, error) {
	return r.queries.ListUsers(ctx)
}

func (r *userRepository) Get(ctx context.Context, id pgtype.UUID) (store.User, error) {
	return r.queries.GetUser(ctx, id)
}

func (r *userRepository) UpdateName(ctx context.Context, params store.UpdateUserNameParams) error {
	return r.queries.UpdateUserName(ctx, params)
}

func (r *userRepository) UpdateType(ctx context.Context, params store.UpdateUserTypeParams) error {
	return r.queries.UpdateUserType(ctx, params)
}

func (r *userRepository) Delete(ctx context.Context, id pgtype.UUID) error {
	return r.queries.DeleteUser(ctx, id)
}
