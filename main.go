package main
import (
	"log"
	"net/http"
	"go.mod/handlers"
	"go.mod/middlewares"
	"go.mod/services"
)
func main() {
	hotelService := services.NewHotelService()
	roomService := services.NewRoomService()
	guestService := services.NewGuestService()
	bookingService := services.NewBookingService()
	
	hotelHandler := &handlers.HotelHandler{Service: hotelService}
	roomHandler := &handlers.RoomHandler{Service: roomService}
	guestHandler := &handlers.GuestHandler{Service: guestService}
	bookingHandler := &handlers.BookingHandler{Service: bookingService}
	
	http.Handle("/hotels", middlewares.LoggingMiddleware(middlewares.AuthMiddleware(hotelHandler)))
	http.Handle("/hotels/", middlewares.LoggingMiddleware(middlewares.AuthMiddleware(hotelHandler)))
	
	http.Handle("/rooms", middlewares.LoggingMiddleware(middlewares.AuthMiddleware(roomHandler)))
	http.Handle("/rooms/", middlewares.LoggingMiddleware(middlewares.AuthMiddleware(roomHandler)))
	
	http.Handle("/guests", middlewares.LoggingMiddleware(middlewares.AuthMiddleware(guestHandler)))
	http.Handle("/guests/", middlewares.LoggingMiddleware(middlewares.AuthMiddleware(guestHandler)))
	
	http.Handle("/bookings", middlewares.LoggingMiddleware(middlewares.AuthMiddleware(bookingHandler)))
	http.Handle("/bookings/", middlewares.LoggingMiddleware(middlewares.AuthMiddleware(bookingHandler)))
	
	port := ":8080"
	log.Printf("Сервер REST API (Handler-Service-Repository) запущено на http://localhost%s", port)
	log.Fatal(http.ListenAndServe(port, nil))
}