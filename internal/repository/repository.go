package repository

import (
	"time"

	"github.com/arkadiuszekprogramista/bookingapp/internal/models"
)

type DatabaseRepo interface {
	AllUsers() bool

	InsertReservation(res models.Reservation) (int, error)
	InsertRoomRestriction(r models.RoomRestriction) error
	SerachAvailabilityByDatesByRoomID(start, end time.Time, roomID int) (bool, error )
	SearchAvailabilityForAllRooms(start, end time.Time) ([]models.Room, error) 
}