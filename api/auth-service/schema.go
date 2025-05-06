package main

import (
	"time"
)

type User struct {
	Name       string    `json:"name"`
	Age        int32     `json:age`
	Mobile     int64     `json:"mobile`
	Email      string    `json:"email`
	Country    string    `json:"country`
	Occupation string    `json:occupation`
	Password   string    `json:password`
	CreatedAt  time.Time `json:created_at"`
}

type Token struct {
	Token     string    `json:"token"`
	CreatedAt time.Time `json:"created_at"`
	ExpiresAt time.Time `json:"expires_at"`
	Email     string    `json:"email"`
}
