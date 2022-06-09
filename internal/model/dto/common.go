package dto

type Response struct {
	Data interface{}
	Code int
}

type ErrorResponse struct {
	Message string `json:"message"`
}
