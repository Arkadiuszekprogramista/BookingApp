package dbrepo

import (
	"errors"
	"time"

	"github.com/arkadiuszekprogramista/bookingapp/internal/models"
)

func (m *testDBRepo) AllUsers() bool {
	return true
	
}

//InsertReservation insert a reservation into database
func (m *testDBRepo) InsertReservation(res models.Reservation) (int, error) {
	// if the room id is 2, then fail; otherwise, pass
	if res.RoomID == 2 {
		return 0, errors.New("some error")
	}

	return 1, nil
}

//InsertRoomRestriction insert a room restrictioninto the database
func (m *testDBRepo) InsertRoomRestriction(r models.RoomRestriction) error {
	if r.RoomID == 1000 {
		return errors.New("some error")
	}
	return nil
}

//SerachAvailabilityByDatesByRoomID returns true if availability exists for roomID, and false if no availability
func (m *testDBRepo) SerachAvailabilityByDatesByRoomID(start, end time.Time, roomID int) (bool, error ){
	return false, nil
}


// SearchAvailabilityForAllRooms slice of rooms for availability rooms, if any, for given dat  range
func (m * testDBRepo) SearchAvailabilityForAllRooms(start, end time.Time) ([]models.Room, error) {
	var rooms []models.Room
	return rooms, nil
}

//GetRoomByID gets a room by id
func(m *testDBRepo) GetRoomByID(id int) (models.Room, error) {
	var room models.Room

	if id != 3 && id != 4 {
		return room, errors.New("some error, room id does not exist")
	}
	return room, nil
}