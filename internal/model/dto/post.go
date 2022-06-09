package dto

type PostData struct {
	Parent  int64  `json:"parent"`
	Author  string `json:"author"`
	Message string `json:"message"`
}
