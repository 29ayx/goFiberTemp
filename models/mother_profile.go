package models
import (
    "time"
)
type Profile struct {
    ID             uint      `gorm:"primaryKey" json:"id"`
    UserID         uint      `json:"user_id"`
    ProfileOwnerType string  `json:"profile_owner_type"`
    
    // Flattened personal details
    Email       string `json:"email"`
    Phone       string `json:"phone"`
    DateOfBirth string `json:"date_of_birth"`
    BloodType   string `json:"blood_type"`

    // Address
    Street     string `json:"street"`
    City       string `json:"city"`
    State      string `json:"state"`
    Country    string `json:"country"`

    // Allergens
    Penicillin  string `json:"penicillin"`
    Eggs        string `json:"eggs"`
    DustMites   string `json:"dust_mites"`

    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}


