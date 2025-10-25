package models

type User struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Age      int    `json:"age"`
	Gender   string `json:"gender"`
	Email    string `json:"email"`
	Password string `json:"password,omitempty"`
}
