package services

import (
	"log"
	"sync"
	"go.mod/models"
	"go.mod/repositories"
	"github.com/google/uuid"
)

const hotelDataPath = "repositories/data/hotels.json"

type HotelService interface {
	GetAll() ([]models.Hotel, error)
	GetByID(id string) (models.Hotel, bool)
	Create(newHotel models.Hotel) (models.Hotel, error)
	Update(id string, updatedHotel models.Hotel) (models.Hotel, bool, error)
	Delete(id string) (bool, error)
}

type hotelServiceImpl struct {
	hotels []models.Hotel
	mutex  *sync.Mutex
}

func NewHotelService() HotelService {
	service := &hotelServiceImpl{
		mutex: &sync.Mutex{},
	}
	
	err := service.loadHotels()
	if err != nil {
		log.Fatalf("Fatal Error: Could not load initial hotel data: %v", err)
	}

	return service
}


func (s *hotelServiceImpl) loadHotels() error {
	loadedHotels, err := repositories.LoadFromFile[models.Hotel](hotelDataPath)
	if err != nil {
		return err
	}
	s.hotels = loadedHotels
	return nil
}


func (s *hotelServiceImpl) saveHotels() error {
	saveableHotels := make([]models.Saveable, len(s.hotels))
	for i, h := range s.hotels {
		saveableHotels[i] = h
	}
	return repositories.SaveToFile(saveableHotels, hotelDataPath)
}




func (s *hotelServiceImpl) GetAll() ([]models.Hotel, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	return s.hotels, nil
}

func (s *hotelServiceImpl) GetByID(id string) (models.Hotel, bool) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	for _, hotel := range s.hotels {
		if hotel.ID == id {
			return hotel, true
		}
	}
	return models.Hotel{}, false
}

func (s *hotelServiceImpl) Create(newHotel models.Hotel) (models.Hotel, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	newHotel.ID = uuid.New().String()
	s.hotels = append(s.hotels, newHotel)

	if err := s.saveHotels(); err != nil {
		return models.Hotel{}, err
	}
	return newHotel, nil
}

func (s *hotelServiceImpl) Update(id string, updatedHotel models.Hotel) (models.Hotel, bool, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	for i, hotel := range s.hotels {
		if hotel.ID == id {
			updatedHotel.ID = id 
			s.hotels[i] = updatedHotel
			
			if err := s.saveHotels(); err != nil {
				return models.Hotel{}, true, err
			}
			return updatedHotel, true, nil
		}
	}
	return models.Hotel{}, false, nil
}

func (s *hotelServiceImpl) Delete(id string) (bool, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	index := -1
	for i, hotel := range s.hotels {
		if hotel.ID == id {
			index = i
			break
		}
	}

	if index == -1 {
		return false, nil 
	}

	s.hotels = append(s.hotels[:index], s.hotels[index+1:]...)

	if err := s.saveHotels(); err != nil {
		return true, err
	}
	return true, nil
}