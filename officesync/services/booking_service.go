package services

import (
	"errors"
	"officesync/models"
	"time"

	"gorm.io/gorm"
)

type BookingService struct {
	DB *gorm.DB
}

func NewBookingService(db *gorm.DB) *BookingService {
	return &BookingService{DB: db}
}

func (s *BookingService) CreateBooking(userID, resourceID uint, start, end time.Time) (*models.Booking, error) {
	if start.Before(time.Now()) {
		return nil, errors.New("cannot book a resource in the past")
	}
	if end.Before(start) || end.Equal(start) {
		return nil, errors.New("end time must be after start time")
	}

	var resource models.Resource
	if err := s.DB.First(&resource, resourceID).Error; err != nil {
		return nil, errors.New("resource not found")
	}

	var overlappingBookings int64
	s.DB.Model(&models.Booking{}).
		Where("resource_id = ? AND status = ?", resourceID, "active").
		Where("start_time < ? AND end_time > ?", end, start).
		Count(&overlappingBookings)

	if overlappingBookings > 0 {
		return nil, errors.New("resource is already booked during this time")
	}

	booking := models.Booking{
		UserID:     userID,
		ResourceID: resourceID,
		StartTime:  start,
		EndTime:    end,
		Status:     "active",
	}

	if err := s.DB.Create(&booking).Error; err != nil {
		return nil, errors.New("failed to create booking")
	}

	return &booking, nil
}

func (s *BookingService) GetUserBookings(userID uint) ([]models.Booking, error) {
	var bookings []models.Booking

	err := s.DB.Preload("Resource").
		Where("user_id = ? AND status = ?", userID, "active").
		Order("start_time asc").
		Find(&bookings).Error

	return bookings, err
}
