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

func (m *postgresDBRepo) InsertReservation(res models.Reservation) error {
	// to cancel transaction if it's taking too long (user loose connection or so)
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `insert into reservation(first_name, last_name, email, phone, start_date, 
                        end_date, room_id, created_at, updated_at)
                        values ($1, $2, $3, $4, $5, $6, $7, $8. $9)`

	//we use context here, so transaction can be canceled
	_, err := m.DB.ExecContext(ctx, stmt,
		res.FirstName,
		res.LastName,
		res.Email,
		res.Phone,
		res.StartDate,
		res.EndDate,
		res.RoomID,
		time.Now(),
		time.Now(),
	)

	if err != nil {
		return err
	}

	return nil
}
