package services

import (
	"log"
	"sync"
	"go.mod/models"
	"go.mod/repositories"
	"github.com/google/uuid"
)

const bookingDataPath = "repositories/data/booking.json"


type BookingService interface {
	GetAll() ([]models.Booking, error)
	GetByID(id string) (models.Booking, bool)
	Create(newBooking models.Booking) (models.Booking, error)
	Update(id string, updatedBooking models.Booking) (models.Booking, bool, error)
	Delete(id string) (bool, error)
}

type bookingServiceImpl struct {
	bookings []models.Booking
	mutex    *sync.Mutex
}

func NewBookingService() BookingService {
	service := &bookingServiceImpl{
		mutex: &sync.Mutex{},
	}
	err := service.loadBookings()
	if err != nil {
		log.Fatalf("Fatal Error: Could not load initial booking data: %v", err)
	}
	return service
}

func (s *bookingServiceImpl) loadBookings() error {
	loadedBookings, err := repositories.LoadFromFile[models.Booking](bookingDataPath)
	if err != nil {
		return err
	}
	s.bookings = loadedBookings
	return nil
}

func (s *bookingServiceImpl) saveBookings() error {
	saveableBookings := make([]models.Saveable, len(s.bookings))
	for i, b := range s.bookings {
		saveableBookings[i] = b
	}
	return repositories.SaveToFile(saveableBookings, bookingDataPath)
}

func (s *bookingServiceImpl) GetAll() ([]models.Booking, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	return s.bookings, nil
}

func (s *bookingServiceImpl) GetByID(id string) (models.Booking, bool) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	for _, booking := range s.bookings {
		if booking.ID == id {
			return booking, true
		}
	}
	return models.Booking{}, false
}

func (s *bookingServiceImpl) Create(newBooking models.Booking) (models.Booking, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	newBooking.ID = uuid.New().String()
	s.bookings = append(s.bookings, newBooking)
	if err := s.saveBookings(); err != nil {
		return models.Booking{}, err
	}
	return newBooking, nil
}

func (s *bookingServiceImpl) Update(id string, updatedBooking models.Booking) (models.Booking, bool, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	for i, booking := range s.bookings {
		if booking.ID == id {
			updatedBooking.ID = id
			s.bookings[i] = updatedBooking
			if err := s.saveBookings(); err != nil {
				return models.Booking{}, true, err
			}
			return updatedBooking, true, nil
		}
	}
	return models.Booking{}, false, nil
}

func (s *bookingServiceImpl) Delete(id string) (bool, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	index := -1
	for i, booking := range s.bookings {
		if booking.ID == id {
			index = i
			break
		}
	}
	if index == -1 {
		return false, nil
	}
	s.bookings = append(s.bookings[:index], s.bookings[index+1:]...)
	if err := s.saveBookings(); err != nil {
		return true, err
	}
	return true, nil
}