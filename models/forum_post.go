package models

import (
    "time"
)

// ForumPost struct linked to the User via Email
type ForumPost struct {
    ID          uint      `gorm:"primaryKey" json:"id"`
    Title       string    `json:"title"`
    Content     string    `json:"content"`
    Category    string    `json:"category"`
    Type        string    `json:"type"`  // e.g., "question", "discussion", etc.
    Email       string    `json:"email"` // Foreign key to the User's email
    Likes       int       `json:"likes"`
    Replies     int       `json:"replies"`
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`
}
