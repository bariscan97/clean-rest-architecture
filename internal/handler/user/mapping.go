package user

import (
	"github.com/bariscan97/clean-rest-architecture/internal/domains"
)

func CreateReqToDomain(u RegisterUserReq) *domains.User {
	return &domains.User{
		UserName: u.UserName,
		Email: u.Email,
		Password: u.Password,
	}
}

func toUserRes(u *domains.User) UserRes {
	return UserRes{
		ID: u.ID,
		UserName: u.UserName,
		ImgUrl: u.ImgUrl,
		Email: u.Email,
	}
}

func ListUserRes(users []*domains.User) []UserRes {
	var ListUsers []UserRes
	
	for _, user := range users {
		ListUsers = append(ListUsers, toUserRes(user))
			
	}
	
	return ListUsers
}
