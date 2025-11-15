package repositories

import (
	"go.mod/models"
	"gorm.io/gorm"
)

type HotelRepository interface {
	GetAll() ([]models.Hotel, error)
	GetByID(id uint) (models.Hotel, error)
	Create(hotel *models.Hotel) error
	Update(hotel *models.Hotel) error
	Delete(id uint) error
}

type hotelRepository struct {
	db *gorm.DB
}

func NewHotelRepository(db *gorm.DB) HotelRepository {
	return &hotelRepository{db: db}
}

func (r *hotelRepository) GetAll() ([]models.Hotel, error) {
	var hotels []models.Hotel
	err := r.db.Find(&hotels).Error
	return hotels, err
}

func (r *hotelRepository) GetByID(id uint) (models.Hotel, error) {
	var hotel models.Hotel
	err := r.db.First(&hotel, id).Error
	return hotel, err
}

func (r *hotelRepository) Create(hotel *models.Hotel) error {
	return r.db.Create(hotel).Error
}

func (r *hotelRepository) Update(hotel *models.Hotel) error {
	return r.db.Save(hotel).Error
}

func (r *hotelRepository) Delete(id uint) error {
	return r.db.Delete(&models.Hotel{}, id).Error
}
