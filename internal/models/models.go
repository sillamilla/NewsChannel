package models

import "time"

type User struct {
	Image string `json:"image"`
	Name  string `json:"name"`
	ID    string `json:"id"`
}

type Channel struct {
	Image     string `json:"image"`
	Name      string `json:"name"`
	ID        string `json:"id"`
	Followers []User `json:"followers"`
}

type PieceOfNews struct {
	Image   string    `json:"image"`
	Title   string    `json:"title"`
	Content string    `json:"content"`
	Date    time.Time `json:"date"`
}

type Notification struct {
	From Channel       `json:"from"`
	News []PieceOfNews `json:"news"`
}
