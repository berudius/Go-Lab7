package services

import (
	"log"
	"sync"
	"go.mod/models"
	"go.mod/repositories"
	"github.com/google/uuid"
)

const roomDataPath = "repositories/data/rooms.json"

type RoomService interface {
	GetAll() ([]models.Room, error)
	GetByID(id string) (models.Room, bool)
	Create(newRoom models.Room) (models.Room, error)
	Update(id string, updatedRoom models.Room) (models.Room, bool, error)
	Delete(id string) (bool, error)
}

type roomServiceImpl struct {
	rooms []models.Room
	mutex *sync.Mutex
}

func NewRoomService() RoomService {
	service := &roomServiceImpl{
		mutex: &sync.Mutex{},
	}
	err := service.loadRooms()
	if err != nil {
		log.Fatalf("Fatal Error: Could not load initial room data: %v", err)
	}
	return service
}

func (s *roomServiceImpl) loadRooms() error {
	loadedRooms, err := repositories.LoadFromFile[models.Room](roomDataPath)
	if err != nil {
		return err
	}
	s.rooms = loadedRooms
	return nil
}

func (s *roomServiceImpl) saveRooms() error {
	saveableRooms := make([]models.Saveable, len(s.rooms))
	for i, r := range s.rooms {
		saveableRooms[i] = r
	}
	return repositories.SaveToFile(saveableRooms, roomDataPath)
}

func (s *roomServiceImpl) GetAll() ([]models.Room, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	return s.rooms, nil
}

func (s *roomServiceImpl) GetByID(id string) (models.Room, bool) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	for _, room := range s.rooms {
		if room.ID == id {
			return room, true
		}
	}
	return models.Room{}, false
}

func (s *roomServiceImpl) Create(newRoom models.Room) (models.Room, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	newRoom.ID = uuid.New().String()
	s.rooms = append(s.rooms, newRoom)
	if err := s.saveRooms(); err != nil {
		return models.Room{}, err
	}
	return newRoom, nil
}

func (s *roomServiceImpl) Update(id string, updatedRoom models.Room) (models.Room, bool, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	for i, room := range s.rooms {
		if room.ID == id {
			updatedRoom.ID = id
			s.rooms[i] = updatedRoom
			if err := s.saveRooms(); err != nil {
				return models.Room{}, true, err
			}
			return updatedRoom, true, nil
		}
	}
	return models.Room{}, false, nil
}

func (s *roomServiceImpl) Delete(id string) (bool, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	index := -1
	for i, room := range s.rooms {
		if room.ID == id {
			index = i
			break
		}
	}
	if index == -1 {
		return false, nil
	}
	s.rooms = append(s.rooms[:index], s.rooms[index+1:]...)
	if err := s.saveRooms(); err != nil {
		return true, err
	}
	return true, nil
}