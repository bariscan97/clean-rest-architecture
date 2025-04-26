package post

import (
	"context"
	"fmt"
	"errors"
	"github.com/bariscan97/clean-rest-architecture/internal/domains"
	"github.com/bariscan97/clean-rest-architecture/internal/utils"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5"
)

type IPostRepository interface {
	ListPosts(ctx context.Context, userID *uuid.UUID, parentID *uuid.UUID, page int, limit int) ([]*domains.PostManyToMany, error)
	CreatePost(ctx context.Context, parentID *uuid.UUID, userID uuid.UUID, post *domains.Post) (*domains.Post, error)
	DeletePostByID(ctx context.Context, userID uuid.UUID, postID uuid.UUID) error
    UpdatePost(ctx context.Context, postID uuid.UUID, userID uuid.UUID, fields map[string]interface{}) error
}

type postRepository struct {
	pool *pgxpool.Pool
}

func NewUserRepository(pool *pgxpool.Pool) IPostRepository {
	return &postRepository{pool: pool}
}

func (r *postRepository) ListPosts(
	ctx context.Context,
	userID *uuid.UUID,
	parentID *uuid.UUID,
	page int,
	limit int,
) ([]*domains.PostManyToMany, error) {

	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	offset := (page - 1) * limit

	var (
		condition string
		params     []any
	)

	index := 1

	if userID != nil {
		condition = fmt.Sprintf(" p.user_id = $%d", index)
		params = append(params, *userID)
		index++
	}else if parentID != nil {
		condition = fmt.Sprintf(" p.parent_id = $%d", index)
		params = append(params, *parentID)
		index++
	}

	whereClause := ""
	if condition != "" {
		whereClause = "WHERE " + condition
	}

	query := fmt.Sprintf(`
        SELECT p.id, p.parent_id, p.user_id, u.username, u.img_url, 
               p.title, p.content, p.update_at, p.create_at
        FROM posts AS p
        LEFT JOIN users AS u 
        ON u.id = p.user_id
        %s
        ORDER BY p.create_at DESC
        LIMIT $%d OFFSET $%d
    `, whereClause, index, index+1)

	params = append(params, limit, offset)

	rows, err := r.pool.Query(ctx, query, params...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []*domains.PostManyToMany

	for rows.Next() {
		var p domains.PostManyToMany
		if err := rows.Scan(
			&p.ID,
			&p.ParentID,
			&p.UserID,
			&p.UserName,
			&p.UserImg,
			&p.Title,
			&p.Content,
			&p.UpdateAt,
			&p.CreateAt,
		); err != nil {
			return nil, err
		}
		posts = append(posts, &p)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return posts, nil
}

func (r *postRepository) GetUserPostsById(ctx context.Context, postID uuid.UUID) (*domains.Post, error) {
	query := `
		SELECT id, parent_id, user_id, title, content, create_at, update_at
		FROM posts
		WHERE id = $1
	`

	var post domains.Post
	if err := r.pool.QueryRow(ctx, query, postID).Scan(
		&post.ID,
		&post.ParentID,
		&post.UserID,
		&post.Title,
		&post.Content,
		&post.CreateAt,
		&post.UpdateAt,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("post not found: %w", err)
		}
		return nil, fmt.Errorf("failed to get post by ID: %w", err)
	}

	return &post, nil
}

func (r *postRepository) UpdatePost(ctx context.Context, userID uuid.UUID, postID uuid.UUID, fields map[string]interface{}) error {
	query, parameters := utils.BuildUpdateQueryMap("users", fields, map[string]interface{}{
		"id": postID,
		"user_id": userID,
	})

	result, err := r.pool.Exec(ctx, query, parameters...)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return fmt.Errorf("no rows updated for postID: %s", postID)
	}

	return nil
}

func (r *postRepository) CreatePost(ctx context.Context, parentID *uuid.UUID, userID uuid.UUID, post *domains.Post) (*domains.Post, error) {

	query := `
        INSERT INTO posts(parent_id, user_id, title, content)
        VALUES ($1, $2, $3, $4)
		RETURNING id, parent_id, user_id, title, content, update_at, create_at
    `
	var p domains.Post

	if err := r.pool.QueryRow(ctx, query, post.ParentID, post.Title, post.Content).Scan(
		&p.ID,
		&p.ParentID,
		&p.UserID,
		&p.Title,
		&p.Content,
		&p.UpdateAt,
		&p.CreateAt,
	); err != nil {
		return nil, err
	}

	return &p, nil
}

func (r *postRepository) DeletePostByID(ctx context.Context, userID uuid.UUID, postID uuid.UUID) error {
	query := `DELETE FROM posts WHERE id = $1 and user_id = $2`
	_, err := r.pool.Exec(ctx, query, postID, userID)
	if err != nil {
		return fmt.Errorf("failed to delete post with id %s: %w", postID, err)
	}
	return nil
}
