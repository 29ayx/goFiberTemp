package database

import (
    "gorm.io/driver/postgres"
    "gorm.io/gorm"
	"myapp/models"
    "log"
)

var DB *gorm.DB

func Connect() {
    dsn := "host=localhost user=myuser password=mypassword dbname=sfours port=5432 sslmode=disable"
    var err error
    DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
    if err != nil {
        log.Fatal("Failed to connect to the database", err)
    }

    // Auto migrate database models
    DB.AutoMigrate(&models.User{}, &models.Car{}, &models.Profile{}, &models.ForumPost{}, &models.DailyEntry{}, &models.Doctor{}, &models.AdminPost{}, &models.AppointmentForm{})
}
