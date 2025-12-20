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
