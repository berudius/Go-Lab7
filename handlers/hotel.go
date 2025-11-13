package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"go.mod/models"
	"go.mod/services"
)


type HotelHandler struct {
	Service services.HotelService
}

func (h *HotelHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	
	pathSegments := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	var id string
	if len(pathSegments) == 2 && pathSegments[0] == "hotels" {
		id = pathSegments[1]
	}

	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case http.MethodGet:
		if id != "" {
			h.getHotelByID(w, r, id)
		} else {
			h.getAllHotels(w, r)
		}
	case http.MethodPost:
		h.createHotel(w, r)
	case http.MethodPut:
		if id != "" {
			h.updateHotel(w, r, id)
		} else {
			http.Error(w, "ID required for update", http.StatusBadRequest)
		}
	case http.MethodDelete:
		if id != "" {
			h.deleteHotel(w, r, id)
		} else {
			http.Error(w, "ID required for delete", http.StatusBadRequest)
		}
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *HotelHandler) getAllHotels(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	name := query.Get("name")
	roomType := query.Get("room_type")

	hotels, err := h.Service.GetAll()
	if err != nil {
		http.Error(w, "Internal server error reading data", http.StatusInternalServerError)
		return
	}

	if name != "" {
		var filtered []models.Hotel
		for _, hotel := range hotels {
			if strings.Contains(strings.ToLower(hotel.Name), strings.ToLower(name)) {
				filtered = append(filtered, hotel)
			}
		}
		hotels = filtered
	}

	if roomType != "" {
		var filtered []models.Hotel
		for _, hotel := range hotels {
			for _, room := range hotel.Rooms {
				if strings.EqualFold(room.RoomType, roomType) {
					filtered = append(filtered, hotel)
					break 
				}
			}
		}
		hotels = filtered
	}

	json.NewEncoder(w).Encode(hotels)
}

func (h *HotelHandler) getHotelByID(w http.ResponseWriter, r *http.Request, id string) {
	hotel, found := h.Service.GetByID(id)
	if !found {
		http.Error(w, "Hotel not found", http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(hotel)
}

func (h *HotelHandler) createHotel(w http.ResponseWriter, r *http.Request) {
	var newHotel models.Hotel
	if err := json.NewDecoder(r.Body).Decode(&newHotel); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	createdHotel, err := h.Service.Create(newHotel)
	if err != nil {
		log.Printf("Error creating hotel: %v", err)
		http.Error(w, "Server error during creation", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(createdHotel)
}

func (h *HotelHandler) updateHotel(w http.ResponseWriter, r *http.Request, id string) {
	var updatedHotel models.Hotel
	if err := json.NewDecoder(r.Body).Decode(&updatedHotel); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	resultHotel, found, err := h.Service.Update(id, updatedHotel)
	if err != nil {
		log.Printf("Error updating hotel: %v", err)
		http.Error(w, "Server error during update", http.StatusInternalServerError)
		return
	}
	if !found {
		http.Error(w, "Hotel not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(resultHotel)
}

func (h *HotelHandler) deleteHotel(w http.ResponseWriter, r *http.Request, id string) {
	found, err := h.Service.Delete(id)
	if err != nil {
		log.Printf("Error deleting hotel: %v", err)
		http.Error(w, "Server error during deletion", http.StatusInternalServerError)
		return
	}
	if !found {
		http.Error(w, "Hotel not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
