package models

import "time"

type Chat struct {
	ID        string    `json:"id" bson:"_id"`
	Title     string    `json:"title" bson:"title"`
	Messages  []Message `json:"messages" bson:"messages"`
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time `json:"updated_at" bson:"updated_at"`
}

type Message struct {
	ID     string    `json:"id" bson:"_id"`
	Text   string    `json:"text" bson:"text"`
	Role   string    `json:"role" bson:"role"`
	AI     string    `json:"ai,omitempty" bson:"ai,omitempty"`
	SentAt time.Time `json:"sent_at" bson:"sent_at"`
}
