package entity

import "time"

type User struct {
	ID        string    `json:"id"`
	OauthID   string    `json:"oauth_id"`
	Email     string    `json:"email"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	BirthDate time.Time `json:"birth_date"`
	Avatar    string    `json:"avatar"`
	Puid      string    `json:"puid"`
}
