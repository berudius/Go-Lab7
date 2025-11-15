package main

import (
	"log"
	"net/http"

	"go.mod/handlers"
	"go.mod/middlewares"
	"go.mod/repositories"
	"go.mod/services"
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


	http.Handle("/hotels", route(hotelHandler))
	http.Handle("/hotels/", route(hotelHandler))

	http.Handle("/rooms", route(roomHandler))
	http.Handle("/rooms/", route(roomHandler))

	http.Handle("/guests", route(guestHandler))
	http.Handle("/guests/", route(guestHandler))

	http.Handle("/bookings", route(bookingHandler))
	http.Handle("/bookings/", route(bookingHandler))


	port := ":8080"
	log.Printf("Сервер REST API запущено на http://localhost%s", port)
	log.Fatal(http.ListenAndServe(port, nil))
}

func route(h http.Handler) http.Handler {
	return middlewares.LoggingMiddleware(
		middlewares.AuthMiddleware(h),
	)
}
