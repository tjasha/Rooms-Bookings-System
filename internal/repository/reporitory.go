package repository

import "github.com/tjasha/Rooms-Bookings-System/internal/models"

// listing all functions that we need in handlers
type DatabaseRepo interface {
	AllUsers() bool

	InsertReservation(res models.Reservation) error
}
