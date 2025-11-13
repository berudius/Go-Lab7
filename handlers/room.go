package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"

	"go.mod/models"
	"go.mod/services"
)

type RoomHandler struct {
	Service services.RoomService
}

func (h *RoomHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	pathSegments := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	var id string
	if len(pathSegments) == 2 && pathSegments[0] == "rooms" {
		id = pathSegments[1]
	}

	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case http.MethodGet:
		if id != "" {
			h.getRoomByID(w, r, id)
		} else {
			h.getAllRooms(w, r)
		}
	case http.MethodPost:
		h.createRoom(w, r)
	case http.MethodPut:
		if id != "" {
			h.updateRoom(w, r, id)
		} else {
			http.Error(w, "ID required for update", http.StatusBadRequest)
		}
	case http.MethodDelete:
		if id != "" {
			h.deleteRoom(w, r, id)
		} else {
			http.Error(w, "ID required for delete", http.StatusBadRequest)
		}
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *RoomHandler) getAllRooms(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	roomType := query.Get("room_type")
	minPriceStr := query.Get("min_price")
	maxPriceStr := query.Get("max_price")

	rooms, err := h.Service.GetAll()
	if err != nil {
		http.Error(w, "Internal server error reading data", http.StatusInternalServerError)
		return
	}

	// Фільтрація за типом кімнати
	if roomType != "" {
		var filtered []models.Room
		for _, room := range rooms {
			if strings.EqualFold(room.RoomType, roomType) {
				filtered = append(filtered, room)
			}
		}
		rooms = filtered
	}

	// Фільтрація за мінімальною ціною
	if minPriceStr != "" {
		minPrice, err := strconv.ParseFloat(minPriceStr, 32)
		if err != nil {
			http.Error(w, "Invalid min_price format", http.StatusBadRequest)
			return
		}
		priceValue := float32(minPrice)

		var filtered []models.Room
		for _, room := range rooms {
			if room.Price >= priceValue {
				filtered = append(filtered, room)
			}
		}
		rooms = filtered
	}

	// Фільтрація за максимальною ціною
	if maxPriceStr != "" {
		maxPrice, err := strconv.ParseFloat(maxPriceStr, 32)
		if err != nil {
			http.Error(w, "Invalid max_price format", http.StatusBadRequest)
			return
		}
		priceValue := float32(maxPrice)

		var filtered []models.Room
		for _, room := range rooms {
			if room.Price <= priceValue {
				filtered = append(filtered, room)
			}
		}
		rooms = filtered
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(rooms)
}

func (h *RoomHandler) getRoomByID(w http.ResponseWriter, r *http.Request, id string) {
	room, found := h.Service.GetByID(id)
	if !found {
		http.Error(w, "Room not found", http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(room)
}

func (h *RoomHandler) createRoom(w http.ResponseWriter, r *http.Request) {
	var newRoom models.Room
	if err := json.NewDecoder(r.Body).Decode(&newRoom); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	createdRoom, err := h.Service.Create(newRoom)
	if err != nil {
		log.Printf("Error creating room: %v", err)
		http.Error(w, "Server error during creation", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(createdRoom)
}

func (h *RoomHandler) updateRoom(w http.ResponseWriter, r *http.Request, id string) {
	var updatedRoom models.Room
	if err := json.NewDecoder(r.Body).Decode(&updatedRoom); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	resultRoom, found, err := h.Service.Update(id, updatedRoom)
	if err != nil {
		log.Printf("Error updating room: %v", err)
		http.Error(w, "Server error during update", http.StatusInternalServerError)
		return
	}
	if !found {
		http.Error(w, "Room not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(resultRoom)
}

func (h *RoomHandler) deleteRoom(w http.ResponseWriter, r *http.Request, id string) {
	found, err := h.Service.Delete(id)
	if err != nil {
		log.Printf("Error deleting room: %v", err)
		http.Error(w, "Server error during deletion", http.StatusInternalServerError)
		return
	}
	if !found {
		http.Error(w, "Room not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
