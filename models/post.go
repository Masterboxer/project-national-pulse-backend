package models

import "time"

type Post struct {
	ID         int       `json:"id"`
	UserID     int       `json:"user_id"`
	TemplateID int       `json:"templateId"`
	Text       string    `json:"text"`
	PhotoPath  *string   `json:"photoPath,omitempty"`
	CreatedAt  time.Time `json:"created_at"`
}

type PostWithUser struct {
	ID          int    `json:"id"`
	UserID      int    `json:"user_id"`
	TemplateID  int    `json:"template_id"`
	Text        string `json:"text"`
	PhotoPath   string `json:"photo_path"`
	CreatedAt   string `json:"created_at"`
	Username    string `json:"username"`
	DisplayName string `json:"display_name"`
}
