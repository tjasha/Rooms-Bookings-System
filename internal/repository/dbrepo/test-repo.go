package dbrepo

// this is a copy of postgres.go. We need all functions defined in the interface
// we need to change all functions to have receiver testDBRepo instead of postgres

import (
	"errors"
	"github.com/tjasha/Rooms-Bookings-System/internal/models"
	"time"
)

// functions for usage in handlers - should be add to postgresDBRepo interface
func (m *testDBRepo) AllUsers() bool {
	return true
}

func (m *testDBRepo) InsertReservation(res models.Reservation) (int, error) {
	return 1, nil
}

// InsertRoomRestriction Inserts room restrictions into the database
func (m *testDBRepo) InsertRoomRestriction(r models.RoomRestriction) error {
	return nil
}

// SearchAvailabilityForDatesByRoomID returns true if availability exist and false if it doesn't got given room
func (m *testDBRepo) SearchAvailabilityForDatesByRoomID(start, end time.Time, roomID int) (bool, error) {
	return false, nil
}

// SearchAvailabilityForAllRooms returns a slice of available rooms for gived date range
func (m *testDBRepo) SearchAvailabilityForAllRooms(start, end time.Time) ([]models.Room, error) {
	var rooms []models.Room
	return rooms, nil
}

// GetRoomById is getting room by id
func (m *testDBRepo) GetRoomById(id int) (models.Room, error) {
	var room models.Room
	if id > 2 {
		return room, errors.New("some error")
	}

	return room, nil
}
