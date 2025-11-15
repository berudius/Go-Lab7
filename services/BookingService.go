package services

import (
	"go.mod/models"
	"go.mod/repositories"
)

type BookingService interface {
	GetAll() ([]models.Booking, error)
	GetByID(id uint) (models.Booking, error)
	Create(booking *models.Booking) error
	Update(booking *models.Booking) error
	Delete(id uint) error
}

type bookingServiceImpl struct {
	repo repositories.BookingRepository
}

func NewBookingService(repo repositories.BookingRepository) BookingService {
	return &bookingServiceImpl{repo: repo}
}

func (s *bookingServiceImpl) GetAll() ([]models.Booking, error) {
	return s.repo.GetAll()
}

func (s *bookingServiceImpl) GetByID(id uint) (models.Booking, error) {
	return s.repo.GetByID(id)
}

func (s *bookingServiceImpl) Create(booking *models.Booking) error {
	return s.repo.Create(booking)
}

func (s *bookingServiceImpl) Update(booking *models.Booking) error {
	return s.repo.Update(booking)
}

func (s *bookingServiceImpl) Delete(id uint) error {
	return s.repo.Delete(id)
}
