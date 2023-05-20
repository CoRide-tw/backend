package model

import "time"

type User struct {
	Id        int32      `json:"id"`
	Name      string     `json:"name"`
	Email     string     `json:"email"`
	GoogleId  int32      `json:"googleId"`
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt time.Time  `json:"updatedAt"`
	DeletedAt *time.Time `json:"deletedAt"`
}
