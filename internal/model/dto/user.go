package dto

type CreateUserRequest struct {
	Nickname string `path:"nickname"`
	Fullname string `json:"fullname"`
	About    string `json:"about"`
	Email    string `json:"email"`
}
