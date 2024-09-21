package models

import (
    "time"
)

type User struct {
    ID        uint      `gorm:"primaryKey" json:"id"`
    FirstName      string    `json:"firstname"`
    LastName      string    `json:"lastname"`
    Email     string    `gorm:"unique" json:"email"`
    Phone     string    `json:"phone"`
    Password  string    `json:"password"`
    Role      string    `json:"role"` // "doctor" or "mother"
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}