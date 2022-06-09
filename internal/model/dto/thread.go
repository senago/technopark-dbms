package dto

import "time"

type CreateForumThreadRequest struct {
	Forum   string    `path:"slug"`
	Title   string    `json:"title"`
	Author  string    `json:"author"`
	Message string    `json:"message"`
	Slug    string    `json:"slug"`
	Created time.Time `json:"created,omitempty"`
}
