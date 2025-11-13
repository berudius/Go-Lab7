package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"go.mod/models"
	"go.mod/services"
)

type BookingHandler struct {
	Service services.BookingService
}

func (h *BookingHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	pathSegments := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	var id string
	if len(pathSegments) == 2 && pathSegments[0] == "bookings" {
		id = pathSegments[1]
	}

	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case http.MethodGet:
		if id != "" {
			h.getBookingByID(w, r, id)
		} else {
			h.getAllBookings(w, r)
		}
	case http.MethodPost:
		h.createBooking(w, r)
	case http.MethodPut:
		if id != "" {
			h.updateBooking(w, r, id)
		} else {
			http.Error(w, "ID required for update", http.StatusBadRequest)
		}
	case http.MethodDelete:
		if id != "" {
			h.deleteBooking(w, r, id)
		} else {
			http.Error(w, "ID required for delete", http.StatusBadRequest)
		}
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *BookingHandler) getAllBookings(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	guestID := query.Get("guest_id")
	roomType := query.Get("room_type")

	bookings, err := h.Service.GetAll()
	if err != nil {
		http.Error(w, "Internal server error reading data", http.StatusInternalServerError)
		return
	}

	// Фільтрація за guest_id
	if guestID != "" {
		var filtered []models.Booking
		for _, booking := range bookings {
			if booking.Guest.ID == guestID {
				filtered = append(filtered, booking)
			}
		}
		bookings = filtered
	}

	// Фільтрація за room_type
	if roomType != "" {
		var filtered []models.Booking
		for _, booking := range bookings {
			for _, room := range booking.BookedRooms {
				if strings.EqualFold(room.RoomType, roomType) {
					filtered = append(filtered, booking)
					break // додаємо бронювання лише один раз
				}
			}
		}
		bookings = filtered
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(bookings)
}

func (h *BookingHandler) getBookingByID(w http.ResponseWriter, r *http.Request, id string) {
	booking, found := h.Service.GetByID(id)
	if !found {
		http.Error(w, "Booking not found", http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(booking)
}

func (h *BookingHandler) createBooking(w http.ResponseWriter, r *http.Request) {
	var newBooking models.Booking
	if err := json.NewDecoder(r.Body).Decode(&newBooking); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	createdBooking, err := h.Service.Create(newBooking)
	if err != nil {
		log.Printf("Error creating booking: %v", err)
		http.Error(w, "Server error during creation", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(createdBooking)
}

func (h *BookingHandler) updateBooking(w http.ResponseWriter, r *http.Request, id string) {
	var updatedBooking models.Booking
	if err := json.NewDecoder(r.Body).Decode(&updatedBooking); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	resultBooking, found, err := h.Service.Update(id, updatedBooking)
	if err != nil {
		log.Printf("Error updating booking: %v", err)
		http.Error(w, "Server error during update", http.StatusInternalServerError)
		return
	}
	if !found {
		http.Error(w, "Booking not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(resultBooking)
}

func (h *BookingHandler) deleteBooking(w http.ResponseWriter, r *http.Request, id string) {
	found, err := h.Service.Delete(id)
	if err != nil {
		log.Printf("Error deleting booking: %v", err)
			http.Error(w, "Server error during deletion", http.StatusInternalServerError)
		return
	}
	if !found {
		http.Error(w, "Booking not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
