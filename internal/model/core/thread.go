package core

import "time"

type Thread struct {
	ID      int64     `json:"id"`
	Title   string    `json:"title"`
	Author  string    `json:"author"`
	Forum   string    `json:"forum"`
	Message string    `json:"message"`
	Votes   int64     `json:"votes"`
	Slug    string    `json:"slug"`
	Created time.Time `json:"created"`
}
