package domains

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID       uuid.UUID
	UserName string
	Email    string
	Password string
	ImgUrl   string
	UpdateAt *time.Time
	CreateAt time.Time
}




