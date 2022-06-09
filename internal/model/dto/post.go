package dto

import "github.com/senago/technopark-dbms/internal/model/core"

type PostData struct {
	Parent  int64  `json:"parent"`
	Author  string `json:"author"`
	Message string `json:"message"`
}

type PostDetails struct {
	Post   *core.Post   `json:"post"`
	Author *core.User   `json:"author,omitempty"`
	Thread *core.Thread `json:"thread,omitempty"`
	Forum  *core.Forum  `json:"forum,omitempty"`
}

type GetPostDetailsRequest struct {
	ID      int64  `path:"id"`
	Related string `query:"related"`
}
