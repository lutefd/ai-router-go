package models

import "time"

type User struct {
	ID    string `json:"id" bson:"_id"`
	Name  string `json:"name" bson:"name"`
	Email string `json:"email" bson:"email"`
	Role  string `json:"role" bson:"role"`
}

type UserChat struct {
	ID        string    `json:"id" bson:"id"`
	User      string    `json:"user" bson:"user"`
	ChatTitle string    `json:"chat_title" bson:"chat_title"`
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time `json:"updated_at" bson:"updated_at"`
}
