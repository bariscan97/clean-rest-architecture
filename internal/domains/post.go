package domains

import (
	"time"

	"github.com/google/uuid"
)

type Post struct {
	ID       uuid.UUID
	UserID   uuid.UUID
	ParentID *uuid.UUID
	Title    string
	Content  string
	UpdateAt time.Time
	CreateAt time.Time
}

type PostManyToMany struct {
	ID       uuid.UUID
	UserID   uuid.UUID
	ParentID *uuid.UUID
	UserName string
	Title    string
	Content  string
	UserImg  *string
	UpdateAt time.Time
	CreateAt time.Time
}
