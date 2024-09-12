package models

import (
    "time"
)

type Car struct {
    ID        uint      `gorm:"primaryKey" json:"id"`
    Make      string    `json:"make"`
    Model     string    `json:"model"`
    Year      int       `json:"year"`
    UserID    uint      `json:"user_id"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}
