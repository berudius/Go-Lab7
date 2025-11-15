package main

import (
	"go.mod/handlers"
	"go.mod/repositories"
	"go.mod/services"
	"github.com/gin-gonic/gin"
	"log"
)

func main() {
	repositories.InitDB()

	hotelRepo := repositories.NewHotelRepository(repositories.DB)
	hotelService := services.NewHotelService(hotelRepo)
	hotelHandler := handlers.NewHotelHandler(hotelService)

	roomRepo := repositories.NewRoomRepository(repositories.DB)
	roomService := services.NewRoomService(roomRepo)
	roomHandler := handlers.NewRoomHandler(roomService)

	guestRepo := repositories.NewGuestRepository(repositories.DB)
	guestService := services.NewGuestService(guestRepo)
	guestHandler := handlers.NewGuestHandler(guestService)

	bookingRepo := repositories.NewBookingRepository(repositories.DB)
	bookingService := services.NewBookingService(bookingRepo)
	bookingHandler := handlers.NewBookingHandler(bookingService)

	r := gin.Default()

	hotelRoutes := r.Group("/hotels")
	{
		hotelRoutes.GET("/", hotelHandler.GetAll)
		hotelRoutes.GET("/:id", hotelHandler.GetByID)
		hotelRoutes.POST("/", hotelHandler.Create)
		hotelRoutes.PUT("/:id", hotelHandler.Update)
		hotelRoutes.DELETE("/:id", hotelHandler.Delete)
	}

	roomRoutes := r.Group("/rooms")
	{
		roomRoutes.GET("/", roomHandler.GetAll)
		roomRoutes.GET("/:id", roomHandler.GetByID)
		roomRoutes.POST("/", roomHandler.Create)
		roomRoutes.PUT("/:id", roomHandler.Update)
		roomRoutes.DELETE("/:id", roomHandler.Delete)
	}

	guestRoutes := r.Group("/guests")
	{
		guestRoutes.GET("/", guestHandler.GetAll)
		guestRoutes.GET("/:id", guestHandler.GetByID)
		guestRoutes.POST("/", guestHandler.Create)
		guestRoutes.PUT("/:id", guestHandler.Update)
		guestRoutes.DELETE("/:id", guestHandler.Delete)
	}

	bookingRoutes := r.Group("/bookings")
	{
		bookingRoutes.GET("/", bookingHandler.GetAll)
		bookingRoutes.GET("/:id", bookingHandler.GetByID)
		bookingRoutes.POST("/", bookingHandler.Create)
		bookingRoutes.PUT("/:id", bookingHandler.Update)
		bookingRoutes.DELETE("/:id", bookingHandler.Delete)
	}

	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
