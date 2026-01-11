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
	// if the room id is 2, then fail, otherwise pass
	if res.RoomID == 2 {
		{
			return 0, errors.New("test - insert reservation failed")
		}
	}
	return 1, nil
}

// InsertRoomRestriction Inserts room restrictions into the database
func (m *testDBRepo) InsertRoomRestriction(r models.RoomRestriction) error {
	if r.RoomID == 1000 {
		{
			return errors.New("test - insert room restriction failed")
		}
	}

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

func (m *testDBRepo) GetUserByID(id int) (models.User, error) {
	var u models.User
	return u, nil
}

func (m *testDBRepo) UpdateUser(u models.User) error {
	return nil

}

func (m *testDBRepo) Authenticate(email, testPassword string) (int, string, error) {
	return 1, "", nil
}

func (m *testDBRepo) AllReservations() ([]models.Reservation, error) {
	var reservations []models.Reservation

	return reservations, nil
}

// AllNewReservations returns a slice of new reservations
func (m *testDBRepo) AllNewReservations() ([]models.Reservation, error) {
	var reservations []models.Reservation

	return reservations, nil
}

// GetReservationById returns one reservation by ID
func (m *testDBRepo) GetReservationById(id int) (models.Reservation, error) {
	var res models.Reservation

	return res, nil
}

// UpdateReservation updates a reservation in the database
func (m *testDBRepo) UpdateReservation(r models.Reservation) error {
	return nil
}

// DeleteReservation deletes one reservation by id
func (m *testDBRepo) DeleteReservation(id int) error {
	return nil
}

// UpdateProcessedForReservation updates processed for a reservation by id
func (m *testDBRepo) UpdateProcessedForReservation(id, processed int) error {
	return nil
}

func (m *testDBRepo) AllRooms() ([]models.Room, error) {
	var rooms []models.Room

	return rooms, nil
}
