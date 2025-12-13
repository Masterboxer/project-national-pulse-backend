package models

type User struct {
	ID          int    `json:"id"`
	Username    string `json:"username"`
	DisplayName string `json:"display_name"`
	Age         int    `json:"age,omitempty"`
	Gender      string `json:"gender,omitempty"`
	Email       string `json:"email"`
	Password    string `json:"password,omitempty"`
	CreatedAt   string `json:"created_at,omitempty"`
}
