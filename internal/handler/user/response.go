package user

import (
	"time"

	"github.com/google/uuid"
)

type UserRes struct {
	ID       uuid.UUID `json:"id"`
	UserName string    `json:"username"`
	Email    string    `json:"email"`
	ImgUrl   string    `json:"img_url"`
}

type LoginUserRes struct {
	AccessToken          string `json:"accessToken"`
	User                 UserRes
	AccessTokenExpiresAt time.Time `json:"access_token_expires_at"`
}
