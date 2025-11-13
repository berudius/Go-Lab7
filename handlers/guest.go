package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"go.mod/models"
	"go.mod/services"
)

type GuestHandler struct {
	Service services.GuestService
}

func (h *GuestHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	pathSegments := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	var id string
	if len(pathSegments) == 2 && pathSegments[0] == "guests" {
		id = pathSegments[1]
	}

	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case http.MethodGet:
		if id != "" {
			h.getGuestByID(w, r, id)
		} else {
			h.getAllGuests(w, r)
		}
	case http.MethodPost:
		h.createGuest(w, r)
	case http.MethodPut:
		if id != "" {
			h.updateGuest(w, r, id)
		} else {
			http.Error(w, "ID required for update", http.StatusBadRequest)
		}
	case http.MethodDelete:
		if id != "" {
			h.deleteGuest(w, r, id)
		} else {
			http.Error(w, "ID required for delete", http.StatusBadRequest)
		}
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *GuestHandler) getAllGuests(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	name := query.Get("name")
	mobileNumber := query.Get("mobile_number")

	guests, err := h.Service.GetAll()
	if err != nil {
		http.Error(w, "Internal server error reading data", http.StatusInternalServerError)
		return
	}

	if name != "" {
		var filtered []models.Guest
		for _, guest := range guests {
			if strings.Contains(strings.ToLower(guest.Name), strings.ToLower(name)) {
				filtered = append(filtered, guest)
			}
		}
		guests = filtered
	}

	if mobileNumber != "" {
		var filtered []models.Guest
		for _, guest := range guests {
			if guest.MobileNumber == mobileNumber {
				filtered = append(filtered, guest)
			}
		}
		guests = filtered
	}

	json.NewEncoder(w).Encode(guests)
}

func (h *GuestHandler) getGuestByID(w http.ResponseWriter, r *http.Request, id string) {
	guest, found := h.Service.GetByID(id)
	if !found {
		http.Error(w, "Guest not found", http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(guest)
}

func (h *GuestHandler) createGuest(w http.ResponseWriter, r *http.Request) {
	var newGuest models.Guest
	if err := json.NewDecoder(r.Body).Decode(&newGuest); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	createdGuest, err := h.Service.Create(newGuest)
	if err != nil {
		log.Printf("Error creating guest: %v", err)
		http.Error(w, "Server error during creation", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(createdGuest)
}

func (h *GuestHandler) updateGuest(w http.ResponseWriter, r *http.Request, id string) {
	var updatedGuest models.Guest
	if err := json.NewDecoder(r.Body).Decode(&updatedGuest); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	resultGuest, found, err := h.Service.Update(id, updatedGuest)
	if err != nil {
		log.Printf("Error updating guest: %v", err)
		http.Error(w, "Server error during update", http.StatusInternalServerError)
		return
	}
	if !found {
		http.Error(w, "Guest not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(resultGuest)
}

func (h *GuestHandler) deleteGuest(w http.ResponseWriter, r *http.Request, id string) {
	found, err := h.Service.Delete(id)
	if err != nil {
		log.Printf("Error deleting guest: %v", err)
		http.Error(w, "Server error during deletion", http.StatusInternalServerError)
		return
	}
	if !found {
		http.Error(w, "Guest not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
