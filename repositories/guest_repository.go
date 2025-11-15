package repositories

import (
	"go.mod/models"
	"gorm.io/gorm"
)

type GuestRepository interface {
	GetAll() ([]models.Guest, error)
	GetByID(id uint) (models.Guest, error)
	Create(guest *models.Guest) error
	Update(guest *models.Guest) error
	Delete(id uint) error
}

type guestRepository struct {
	db *gorm.DB
}

func NewGuestRepository(db *gorm.DB) GuestRepository {
	return &guestRepository{db: db}
}

func (r *guestRepository) GetAll() ([]models.Guest, error) {
	var guests []models.Guest
	err := r.db.Find(&guests).Error
	return guests, err
}

func (r *guestRepository) GetByID(id uint) (models.Guest, error) {
	var guest models.Guest
	err := r.db.First(&guest, id).Error
	return guest, err
}

func (r *guestRepository) Create(guest *models.Guest) error {
	return r.db.Create(guest).Error
}

func (r *guestRepository) Update(guest *models.Guest) error {
	return r.db.Save(guest).Error
}

func (r *guestRepository) Delete(id uint) error {
	return r.db.Delete(&models.Guest{}, id).Error
}
