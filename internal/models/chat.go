package models

type Chat struct {
	ID       string    `json:"id"`
	Title    string    `json:"title"`
	Messages []Message `json:"messages"`
}

type Message struct {
	ID   string `json:"id"`
	Text string `json:"text"`
	Role string `json:"role"`
	AI   string `json:"ai,omitempty"`
}
