package services

import (
	"log"
	"sync"
	"go.mod/models"
	"go.mod/repositories"
	"github.com/google/uuid"
)

const guestDataPath = "repositories/data/guests.json"

type GuestService interface {
	GetAll() ([]models.Guest, error)
	GetByID(id string) (models.Guest, bool)
	Create(newGuest models.Guest) (models.Guest, error)
	Update(id string, updatedGuest models.Guest) (models.Guest, bool, error)
	Delete(id string) (bool, error)
}

type guestServiceImpl struct {
	guests []models.Guest
	mutex  *sync.Mutex
}

func NewGuestService() GuestService {
	service := &guestServiceImpl{
		mutex: &sync.Mutex{},
	}
	err := service.loadGuests()
	if err != nil {
		log.Fatalf("Fatal Error: Could not load initial guest data: %v", err)
	}
	return service
}

func (s *guestServiceImpl) loadGuests() error {
	loadedGuests, err := repositories.LoadFromFile[models.Guest](guestDataPath)
	if err != nil {
		return err
	}
	s.guests = loadedGuests
	return nil
}

func (s *guestServiceImpl) saveGuests() error {
	saveableGuests := make([]models.Saveable, len(s.guests))
	for i, g := range s.guests {
		saveableGuests[i] = g
	}
	return repositories.SaveToFile(saveableGuests, guestDataPath)
}

func (s *guestServiceImpl) GetAll() ([]models.Guest, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	return s.guests, nil
}

func (s *guestServiceImpl) GetByID(id string) (models.Guest, bool) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	for _, guest := range s.guests {
		if guest.ID == id {
			return guest, true
		}
	}
	return models.Guest{}, false
}

func (s *guestServiceImpl) Create(newGuest models.Guest) (models.Guest, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	newGuest.ID = uuid.New().String()
	s.guests = append(s.guests, newGuest)
	if err := s.saveGuests(); err != nil {
		return models.Guest{}, err
	}
	return newGuest, nil
}

func (s *guestServiceImpl) Update(id string, updatedGuest models.Guest) (models.Guest, bool, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	for i, guest := range s.guests {
		if guest.ID == id {
			updatedGuest.ID = id
			s.guests[i] = updatedGuest
			if err := s.saveGuests(); err != nil {
				return models.Guest{}, true, err
			}
			return updatedGuest, true, nil
		}
	}
	return models.Guest{}, false, nil
}

func (s *guestServiceImpl) Delete(id string) (bool, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	index := -1
	for i, guest := range s.guests {
		if guest.ID == id {
			index = i
			break
		}
	}
	if index == -1 {
		return false, nil
	}
	s.guests = append(s.guests[:index], s.guests[index+1:]...)
	if err := s.saveGuests(); err != nil {
		return true, err
	}
	return true, nil
}