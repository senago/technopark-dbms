package dto

import "github.com/senago/technopark-dbms/internal/model/core"

type CreateUserRequest struct {
	Nickname string `path:"nickname"`
	Fullname string `json:"fullname"`
	About    string `json:"about"`
	Email    string `json:"email"`
}

type GetUserProfileRequest struct {
	Nickname string `path:"nickname"`
}

type GetUserProfileResponse = core.User
