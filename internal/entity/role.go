package entity

type Role struct {
	ID       int64  `json:"id"`
	RoleName string `json:"role_name"`
	UserID   string `json:"user_id"`
}
