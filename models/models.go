package models

import "gorm.io/gorm"

type Hotel struct {
	gorm.Model
	Name  string `gorm:"unique;not null"`
	Rooms []Room `gorm:"foreignKey:HotelID"`
}

type Room struct {
	gorm.Model
	RoomType   string      `gorm:"not null"`
	Price      float32     `gorm:"not null"`
	Facilities StringSlice `gorm:"type:json"`
	HotelID    uint
}

type Guest struct {
	gorm.Model
	Name         string      `gorm:"not null"`
	MobileNumber string      `gorm:"unique;not null"`
	Preferences  StringSlice `gorm:"type:json"`
}

type Booking struct {
	gorm.Model
	GuestID     uint
	HotelID     uint
	Guest       Guest `gorm:"foreignKey:GuestID"`
	Hotel       Hotel `gorm:"foreignKey:HotelID"`
	BookedRooms []Room `gorm:"many2many:booking_rooms;"`
}
