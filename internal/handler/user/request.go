package user

type RegisterUserReq struct {
	UserName string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginUserReq struct {
	Identifier string `json:"identifier"`
	Password   string `json:"password"`
}

type UpdateUserReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	ImgUrl   string `json:"img_url"`
}
