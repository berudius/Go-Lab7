package repositories

import (
	"go.mod/models"
	"gorm.io/gorm"
)

type RoomRepository interface {
	GetAll() ([]models.Room, error)
	GetByID(id uint) (models.Room, error)
	Create(room *models.Room) error
	Update(room *models.Room) error
	Delete(id uint) error
}

type roomRepository struct {
	db *gorm.DB
}

func NewRoomRepository(db *gorm.DB) RoomRepository {
	return &roomRepository{db: db}
}

func (r *roomRepository) GetAll() ([]models.Room, error) {
	var rooms []models.Room
	err := r.db.Find(&rooms).Error
	return rooms, err
}

func (r *roomRepository) GetByID(id uint) (models.Room, error) {
	var room models.Room
	err := r.db.First(&room, id).Error
	return room, err
}

func (r *roomRepository) Create(room *models.Room) error {
	return r.db.Create(room).Error
}

func (r *roomRepository) Update(room *models.Room) error {
	return r.db.Save(room).Error
}

func (r *roomRepository) Delete(id uint) error {
	return r.db.Delete(&models.Room{}, id).Error
}
