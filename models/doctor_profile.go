package models

import (
	"time"
)

type Doctor struct{

	ID uint `gorm:"primaryKey" json:"id"`
	UserID         uint      `json:"user_id"`
Email       string `json:"email"`
	AccStatus string `json:"account_status"`
	 ProfileOwnerType string  `json:"profile_owner_type"`

	PracticeNo string `json:"practice_no"`

	
	FirstName string `json:"first_name"`
	LastName string `json:"last_name"`

	Specialist string `json:"Specialist"`



	
	

	

	// Address
	Street     string `json:"street"`
	City       string `json:"city"`
	State      string `json:"state"`
	Country    string `json:"country"`


	DateOfBirth time.Time `json:"date_of_birth"`



}

