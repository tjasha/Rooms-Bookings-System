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

	SearchAvailabilityByDates(start, end time.Time, roomID int) (bool, error)
}
