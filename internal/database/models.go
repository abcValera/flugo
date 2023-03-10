// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.17.0

package database

import (
	"time"
)

type Joke struct {
	ID          int32     `json:"id"`
	Author      string    `json:"author"`
	Title       string    `json:"title"`
	Text        string    `json:"text"`
	Explanation string    `json:"explanation"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type User struct {
	ID             int32     `json:"id"`
	Username       string    `json:"username"`
	Email          string    `json:"email"`
	HashedPassword string    `json:"hashed_password"`
	Avatar         string    `json:"avatar"`
	Fullname       string    `json:"fullname"`
	Bio            string    `json:"bio"`
	Status         string    `json:"status"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}
