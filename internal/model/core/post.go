package core

import "time"

type Post struct {
	ID       int64     `json:"id"`
	Parent   int64     `json:"parent"`
	Author   string    `json:"author"`
	Message  string    `json:"message"`
	IsEdited bool      `json:"isEdited"`
	Forum    string    `json:"forum"`
	Thread   int64     `json:"thread"`
	Created  time.Time `json:"created"`
}
