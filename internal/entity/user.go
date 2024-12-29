package entity

type User struct {
	ID        string `json:"-"`
	OauthID   string `json:"-"`
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Avatar    string `json:"avatar"`
}

type UserWithTypes struct {
	User  User
	Types []EventType
}
