package repository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"sipasdos/app/model"
)

// UserRepository adalah kontrak akses data akun pengguna.
type UserRepository interface {
	Create(ctx context.Context, u model.User) (model.User, error)
	FindByEmail(ctx context.Context, email string) (model.User, error)
	FindByID(ctx context.Context, id int) (model.User, error)
}

type userRepository struct {
	pool *pgxpool.Pool
}

func NewUserRepository(pool *pgxpool.Pool) UserRepository {
	return &userRepository{pool: pool}
}

const userColumns = "id, name, email, password_hash, role, created_at"

func scanUser(row pgx.Row) (model.User, error) {
	var u model.User
	err := row.Scan(&u.ID, &u.Name, &u.Email, &u.PasswordHash, &u.Role, &u.CreatedAt)
	return u, err
}

func (r *userRepository) Create(ctx context.Context, u model.User) (model.User, error) {
	created, err := scanUser(r.pool.QueryRow(ctx,
		`INSERT INTO users (name, email, password_hash, role)
		 VALUES ($1, $2, $3, $4)
		 RETURNING `+userColumns,
		u.Name, u.Email, u.PasswordHash, u.Role))
	if err != nil {
		return model.User{}, mapWriteError(err) // email ganda -> ErrDuplicate
	}
	return created, nil
}

func (r *userRepository) FindByEmail(ctx context.Context, email string) (model.User, error) {
	u, err := scanUser(r.pool.QueryRow(ctx,
		"SELECT "+userColumns+" FROM users WHERE email = $1", email))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.User{}, ErrNotFound
		}
		return model.User{}, err
	}
	return u, nil
}

func (r *userRepository) FindByID(ctx context.Context, id int) (model.User, error) {
	u, err := scanUser(r.pool.QueryRow(ctx,
		"SELECT "+userColumns+" FROM users WHERE id = $1", id))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.User{}, ErrNotFound
		}
		return model.User{}, err
	}
	return u, nil
}
