package models
import (
	"time"
)
type AdminPost struct {
	ID          uint   `gorm:"primaryKey" json:"id"`

	// Post details
	Title   string `json:"title"`
	Content string `json:"content"`
	Description string `json:"description"`
	Thumbnail string `json:"thumbnail"`
	PostType string `json:"post_type"`
    CreatedAt time.Time `json:"created_at"`
}
