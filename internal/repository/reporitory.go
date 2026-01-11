package repository

import (
	"github.com/tjasha/Rooms-Bookings-System/internal/models"
	"time"
)

// listing all functions that we need in handlers
type DatabaseRepo interface {
	AllUsers() bool

	InsertReservation(res models.Reservation) (int, error)

	InsertRoomRestriction(res models.RoomRestriction) error

	SearchAvailabilityForDatesByRoomID(start, end time.Time, roomID int) (bool, error)

	SearchAvailabilityForAllRooms(start, end time.Time) ([]models.Room, error)

	GetRoomById(id int) (models.Room, error)

	GetUserByID(id int) (models.User, error)

	UpdateUser(u models.User) error

	Authenticate(email, testPassword string) (int, string, error)

	AllReservations() ([]models.Reservation, error)

	AllNewReservations() ([]models.Reservation, error)
}
