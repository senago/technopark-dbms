package dto

type CreateForumRequest struct {
	Title string `json:"title"`
	User  string `json:"user"`
	Slug  string `json:"slug"`
}

type GetForumBySlugRequest struct {
	Slug string `path:"slug"`
}
