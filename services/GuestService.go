package services

import (
	"go.mod/models"
	"go.mod/repositories"
)

type GuestService interface {
	GetAll() ([]models.Guest, error)
	GetByID(id uint) (models.Guest, error)
	Create(guest *models.Guest) error
	Update(guest *models.Guest) error
	Delete(id uint) error
}

type guestServiceImpl struct {
	repo repositories.GuestRepository
}

func NewGuestService(repo repositories.GuestRepository) GuestService {
	return &guestServiceImpl{repo: repo}
}

func (s *guestServiceImpl) GetAll() ([]models.Guest, error) {
	return s.repo.GetAll()
}

func (s *guestServiceImpl) GetByID(id uint) (models.Guest, error) {
	return s.repo.GetByID(id)
}

func (s *guestServiceImpl) Create(guest *models.Guest) error {
	return s.repo.Create(guest)
}

func (s *guestServiceImpl) Update(guest *models.Guest) error {
	return s.repo.Update(guest)
}

func (s *guestServiceImpl) Delete(id uint) error {
	return s.repo.Delete(id)
}
