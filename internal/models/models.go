package models

import "time"

//do store data from/to database

// make reservation form
type Reservation struct {
	FirstName string
	LastName  string
	Email     string
	Phone     string
}

// by convenntion we change names in database first_name to FirstName
// describes User table -> User model
type User struct {
	ID          int
	FirstName   string
	LastName    string
	Email       string
	Password    string
	AccessLevel string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// Rooms is the room model
type Rooms struct {
	ID        int
	RoomName  string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// Restrictions is the restrictions model
type Restrictions struct {
	ID               int
	RestrictionsName string
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

// Reservations is reservation model
type Reservations struct {
	ID        int
	FirstName string
	LastName  string
	Email     string
	Phone     string
	StartDate time.Time
	EndDate   time.Time
	RoomID    int
	CreatedAt time.Time
	UpdatedAt time.Time
	//we can add additional information
	//here we add information about the Room, as it's connected
	Room Rooms
}

// RoomRestriction is roomrestriction model
type RoomRestrictions struct {
	ID            int
	StartDate     time.Time
	EndDate       time.Time
	RoomID        int
	ReservationID int
	RestrictionID int
	CreatedAt     time.Time
	UpdatedAt     time.Time
	// additional information added
	Room        Rooms
	Reservation Reservations
	Restriction Restrictions
}
