package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/bariscan97/clean-rest-architecture/internal/domains"
	"github.com/bariscan97/clean-rest-architecture/internal/utils"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type IUserRepository interface {
	CreateUser(ctx context.Context, user *domains.User) (*domains.User, error)
	ListUsers(ctx context.Context, page int, limit int) ([]*domains.User, error)
	GetUserByIdentifier(ctx context.Context, identifier string) (*domains.User, error)
	UpdateUserByID(ctx context.Context, userID uuid.UUID, fields map[string]interface{}) error
	DeleteUserByID(ctx context.Context, id uuid.UUID) error
}

type userRepository struct {
	pool *pgxpool.Pool
}

func NewUserRepository(pool *pgxpool.Pool) IUserRepository {
	return &userRepository{pool: pool}
}

func (r *userRepository) CreateUser(ctx context.Context, user *domains.User) (*domains.User, error) {
	query := `
		INSERT INTO users (username, email, password)
		VALUES ($1, $2, $3)
		RETURNING id, username, email;
	`
	row := r.pool.QueryRow(ctx, query, user.Email, user.UserName, user.Password)
	var u domains.User
	err := row.Scan(&u.ID, &u.Email, &u.Password)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *userRepository) ListUsers(ctx context.Context, page int, limit int) ([]*domains.User, error) {

	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	offset := (page - 1) * limit

	query := `
		SELECT id, username, img_url, created_at
		FROM users
		ORDER BY created_at
		LIMIT $1 OFFSET $2
	`
	rows, err := r.pool.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*domains.User
	for rows.Next() {
		var u domains.User
		if err := rows.Scan(&u.ID, &u.UserName, &u.ImgUrl, &u.CreateAt); err != nil {
			return nil, err
		}
		users = append(users, &u)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return users, nil
}

func (r *userRepository) GetUserByIdentifier(ctx context.Context, identifier string) (*domains.User, error) {
	query := `
		SELECT id, username, img_url, email, password
		FROM users
		WHERE email = $1 or username = $1 or id = $1;
	`
	row := r.pool.QueryRow(ctx, query, identifier)

	var u domains.User

	err := row.Scan(&u.ID, &u.UserName, &u.ImgUrl, &u.Email, &u.Password)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("user not found for identifier: %s", identifier)
		}
		return nil, err
	}
	return &u, nil
}

func (r *userRepository) UpdateUserByID(ctx context.Context, userID uuid.UUID, fields map[string]interface{}) error {
	sql, parameters := utils.BuildUpdateQueryMap("users", fields, map[string]interface{}{
		"id": userID,
	})

	result, err := r.pool.Exec(ctx, sql, parameters...)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return fmt.Errorf("no rows updated for userID: %s", userID)
	}

	return nil
}

func (r *userRepository) DeleteUserByID(ctx context.Context, id uuid.UUID) error {
	query := `
		DELETE FROM users
		WHERE id = $1;
	`
	_, err := r.pool.Exec(ctx, query, id)
	return err
}
