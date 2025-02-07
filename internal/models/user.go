package models

type User struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Role  string `json:"role"`
}

type UserChat struct {
	ID        string `json:"id"`
	User      string `json:"user"`
	ChatTitle string `json:"chat_title"`
}
