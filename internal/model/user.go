package model

import "time"

type User struct {
	Id         int32      `json:"id"`
	Name       string     `json:"name"`
	Email      string     `json:"email"`
	GoogleId   string     `json:"googleId"`
	PictureUrl string     `json:"pictureUrl"`
	CreatedAt  time.Time  `json:"createdAt"`
	UpdatedAt  time.Time  `json:"updatedAt"`
	DeletedAt  *time.Time `json:"deletedAt"`
}
