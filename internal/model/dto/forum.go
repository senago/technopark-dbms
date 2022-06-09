package dto

type CreateForumRequest struct {
	Title string `json:"title"`
	User  string `json:"user"`
	Slug  string `json:"slug"`
}

type GetForumBySlugRequest struct {
	Slug string `path:"slug"`
}

type GetForumThreadsRequest struct {
	Slug  string `path:"slug"`
	Limit int64  `query:"limit"`
	Since string `query:"since"`
	Desc  bool   `query:"desc"`
}
