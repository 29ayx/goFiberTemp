package models

import (
    "time"
)

type User struct {
    ID        uint      `gorm:"primaryKey" json:"id"`
    Name      string    `json:"name"`
    Email     string    `gorm:"unique" json:"email"`
    Password  string    `json:"password"`
    Cars      []Car     `gorm:"foreignKey:UserID" json:"cars,omitempty"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}
