package post

import (
	"time"

	"github.com/google/uuid"
)

type CreatePostRes struct {
	ID       uuid.UUID  `json:"id"`
	UserID   uuid.UUID  `json:"user_id"`
	ParentID *uuid.UUID `json:"parent_id,omitempty"`
	Title    string     `json:"title"`
	Content  string     `json:"content"`
	UpdateAt time.Time  `json:"update_at"`
	CreateAt time.Time  `json:"create_at"`
}

type FetchPostRes struct {
	ID       uuid.UUID  `json:"id"`
	UserID   uuid.UUID  `json:"user_id"`
	ParentID *uuid.UUID `json:"parent_id,omitempty"`
	UserName string     `json:"username"`
	Title    string     `json:"title"`
	Content  string     `json:"content"`
	UserImg  *string    `json:"user_img,omitempty"`
	UpdateAt time.Time  `json:"update_at"`
	CreateAt time.Time  `json:"create_at"`
}
