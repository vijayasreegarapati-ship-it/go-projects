package services

import (
	"officesync/models"
	"testing"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	db.AutoMigrate(&models.User{}, &models.Resource{}, &models.Booking{})
	return db
}

func TestDoubleBookingPrevention(t *testing.T) {
	db := setupTestDB()
	service := NewBookingService(db)

	user := models.User{Name: "Test User", Email: "test@example.com"}
	db.Create(&user)

	resource := models.Resource{Name: "Desk 1", ResourceType: "desk"}
	db.Create(&resource)

	now := time.Now()

	start1 := now.Add(24 * time.Hour)
	end1 := now.Add(26 * time.Hour)

	start2 := now.Add(25 * time.Hour)
	end2 := now.Add(27 * time.Hour)

	_, err := service.CreateBooking(user.ID, resource.ID, start1, end1)
	if err != nil {
		t.Fatalf("Expected first booking to succeed, but got error: %v", err)
	}

	_, err = service.CreateBooking(user.ID, resource.ID, start2, end2)
	if err == nil {
		t.Fatalf("Expected second booking to fail due to overlap, but it magically succeeded!")
	}

	expectedErr := "resource is already booked during this time"
	if err.Error() != expectedErr {
		t.Errorf("Expected error message '%s', but got '%v'", expectedErr, err)
	}
}
