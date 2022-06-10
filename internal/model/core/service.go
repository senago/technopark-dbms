package core

type ServiceInfo struct {
	User   int64 `json:"user"`
	Forum  int64 `json:"forum"`
	Thread int64 `json:"thread"`
	Post   int64 `json:"post"`
}
