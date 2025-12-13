package models

type User struct {
	ID          int    `json:"id"`
	Username    string `json:"username"`
	DisplayName string `json:"display_name"`
	Age         int    `json:"age"`
	Gender      string `json:"gender"`
	Email       string `json:"email"`
	Password    string `json:"password,omitempty"`
	CreatedAt   string `json:"created_at"`
}
