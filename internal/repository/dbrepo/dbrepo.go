package dbrepo

import (
	"database/sql"
	"github.com/tjasha/Rooms-Bookings-System/internal/config"
	"github.com/tjasha/Rooms-Bookings-System/internal/repository"
)

// Repository pattern:
// this is allowing me to swap databases very quickly,
// just create new type and new function with needed information

type postgresDBRepo struct {
	App *config.AppConfig
	DB  *sql.DB
}

func NewPostresRepo(conn *sql.DB, a *config.AppConfig) repository.DatabaseRepo {
	return &postgresDBRepo{
		App: a,
		DB:  conn,
	}
}
