package models

type Saveable interface {
	Pass()
}

type Hotel struct {
	ID    string
	Name  string
	Rooms []Room
}

type Room struct {
	ID         string
	RoomType   string
	Price      float32
	Facilities []string
}

type Guest struct {
	ID           string
	Name         string
	MobileNumber string
	Preferences  []string
}

type Booking struct {
	ID          string
	Guest       Guest
	Hotel       Hotel
	BookedRooms []Room
}

func (h Hotel) Pass()   {}
func (r Room) Pass()    {}
func (g Guest) Pass()   {}
func (b Booking) Pass() {}
