package models

import (
    "time"
)

type DailyEntry struct {
    ID          uint      `gorm:"primaryKey" json:"id"`
	UserEmail	string	  `json:user_email`
    DailyRating int       `json:"daily_rating"`      // Rating for the day (e.g., 4)
    EntryDate   string    `json:"entry_date"`        // Entry date (e.g., "18-11-2024")
    Feeling     string    `json:"feeling"`           // How the user is feeling (e.g., "good")
    Gratitudes  string    `json:"gratitudes"`        // Things the user is grateful for (e.g., "seeing a butterfly")
    SelfCare    string    `json:"selfcare"`          // Self-care activity (e.g., "haircare")
    Thoughts    string    `json:"thoughts"`          // Thoughts of the user (e.g., "feeling good thoughts")
    UserID      string    `json:"user_id"`           // Foreign key referencing the user's email (e.g., "test5@test.com")
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`
}