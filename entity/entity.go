package entity

import "time"

type User struct {
	CreatedAt   time.Time `json:"created_at"`
	DisplayName string    `json:"display_name"`
	Email       string    `json:"email"`
}

type UserStore struct {
	Increment int              `json:"increment"`
	List      map[string]*User `json:"list"`
}
