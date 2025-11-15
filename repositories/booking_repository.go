package repositories

import (
	"go.mod/models"
	"gorm.io/gorm"
)

type BookingRepository interface {
	GetAll() ([]models.Booking, error)
	GetByID(id uint) (models.Booking, error)
	Create(booking *models.Booking) error
	Update(booking *models.Booking) error
	Delete(id uint) error
}

type bookingRepository struct {
	db *gorm.DB
}

func NewBookingRepository(db *gorm.DB) BookingRepository {
	return &bookingRepository{db: db}
}

func (r *bookingRepository) GetAll() ([]models.Booking, error) {
	var bookings []models.Booking
	err := r.db.Preload("Guest").Preload("Hotel").Preload("BookedRooms").Find(&bookings).Error
	return bookings, err
}

func (r *bookingRepository) GetByID(id uint) (models.Booking, error) {
	var booking models.Booking
	err := r.db.Preload("Guest").Preload("Hotel").Preload("BookedRooms").First(&booking, id).Error
	return booking, err
}

func (r *bookingRepository) Create(booking *models.Booking) error {
	return r.db.Create(booking).Error
}

func (r *bookingRepository) Update(booking *models.Booking) error {
	return r.db.Save(booking).Error
}

func (r *bookingRepository) Delete(id uint) error {
	return r.db.Delete(&models.Booking{}, id).Error
}
