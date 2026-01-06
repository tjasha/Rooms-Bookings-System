package dbrepo

import (
	"context"
	"github.com/tjasha/Rooms-Bookings-System/internal/models"
	"time"
)

// functions for usage in handlers - should be add to postgresDBRepo interface
func (m *postgresDBRepo) AllUsers() bool {
	return true
}

func (m *postgresDBRepo) InsertReservation(res models.Reservation) (int, error) {
	// to cancel transaction if it's taking too long (user loose connection or so)
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var newID int

	stmt := `insert into reservations(first_name, last_name, email, phone, start_date, 
                        end_date, room_id, created_at, updated_at)
                        values ($1, $2, $3, $4, $5, $6, $7, $8, $9) returning id`
	// we need to know the id under which reservation is inserted, so we use returning id

	//we use context here, so transaction can be canceled
	//_, err := m.DB.ExecContext(ctx, stmt,
	// QueryRow context returns only error and has Scan function, which allows us to save id
	err := m.DB.QueryRowContext(ctx, stmt,
		res.FirstName,
		res.LastName,
		res.Email,
		res.Phone,
		res.StartDate,
		res.EndDate,
		res.RoomID,
		time.Now(),
		time.Now(),
	).Scan(&newID)

	if err != nil {
		return 0, err
	}

	return newID, nil
}

// InsertRoomRestriction Inserts room restrictions into the database
func (m *postgresDBRepo) InsertRoomRestriction(r models.RoomRestriction) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `insert into room_restrictions(start_date, end_date, room_id, reservation_id,
                              restriction_id, created_at, updated_at)
                              values ($1, $2, $3, $4, $5, $6, $7)`

	_, err := m.DB.ExecContext(ctx, stmt,
		r.StartDate,
		r.EndDate,
		r.RoomID,
		r.ReservationID,
		r.RestrictionID,
		time.Now(),
		time.Now(),
	)

	if err != nil {
		return err
	}

	return nil
}

// SearchAvailabilityByDates returns true if availability exist and false if it doesn't got given room
func (m *postgresDBRepo) SearchAvailabilityByDates(start, end time.Time, roomID int) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var numRows int
	query := `
		select 
			count(id)
		from
			room_restrictions
		where 
		    roomId = &1 and
		    $2 < end_date and $3 > start_date`

	row := m.DB.QueryRowContext(ctx, query, start, end, roomID)
	err := row.Scan(&numRows)
	if err != nil {
		return false, err
	}

	if numRows == 0 {
		return true, nil
	}

	return true, nil
}
