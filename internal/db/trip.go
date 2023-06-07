package db

import (
	"context"
	"time"

	. "github.com/CoRide-tw/backend/internal/errors/generated/dberr"
	. "github.com/DenChenn/blunder/pkg/blunder"
	"github.com/jackc/pgx/v5"

	"github.com/CoRide-tw/backend/internal/model"
)

const createTripTableSQL = `
	CREATE TABLE IF NOT EXISTS trips (
		id SERIAL,
		rider_id INT NOT NULL,
		driver_id INT NOT NULL,
		request_id INT NOT NULL,
		route_id INT NOT NULL,
		created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL,
		deleted_at TIMESTAMP
	);
`

func initTripTable() error {
	if _, err := DBClient.pgPool.Exec(context.Background(), createTripTableSQL); err != nil {
		return err
	}
	return nil
}

const listTripByRiderIdSQL = `
	SELECT 
		t.id,
		t.rider_id,
		t.driver_id,
		t.request_id,
		t.route_id,
		u.name,
		u.picture_url,
		u.car_type,
		u.car_plate,
		rout.start_time,
		req.pickup_start_time,
		req.pickup_end_time,
		rout.end_time,
		ST_X(rout.start_location), ST_Y(rout.start_location), 
		ST_X(rout.end_location), ST_Y(rout.end_location), 
		ST_X(req.pickup_location), ST_Y(req.pickup_location),
		ST_X(req.dropoff_location), ST_Y(req.dropoff_location),
		t.created_at,
		t.deleted_at
	FROM trips t
		JOIN users u ON t.driver_id = u.id
		JOIN requests req ON t.request_id = req.id
		JOIN routes rout ON t.route_id = rout.id
	WHERE t.rider_id = $1 AND t.deleted_at IS NULL;
`

type ListTripResp struct {
	Id                    int32      `json:"id"`
	RiderId               int32      `json:"riderId"`
	DriverId              int32      `json:"driverId"`
	RequestId             int32      `json:"requestId"`
	RouteId               int32      `json:"routeId"`
	DriverName            string     `json:"driverName"`
	DriverPictureUrl      string     `json:"driverPictureUrl"`
	DriverCarType         string     `json:"driverCarType"`
	DriverCarPlate        string     `json:"driverCarPlate"`
	RouteStartTime        time.Time  `json:"routeStartTime"`
	PickupStartTime       time.Time  `json:"pickupStartTime"`
	PickupEndTime         time.Time  `json:"pickupEndTime"`
	RouteEndTime          time.Time  `json:"routeEndTime"`
	RouteStartLocationLng float64    `json:"routeStartLocationLng"`
	RouteStartLocationLat float64    `json:"routeStartLocationLat"`
	RouteEndLocationLng   float64    `json:"routeEndLocationLng"`
	RouteEndLocationLat   float64    `json:"routeEndLocationLat"`
	PickupLocationLng     float64    `json:"pickupLocationLng"`
	PickupLocationLat     float64    `json:"pickupLocationLat"`
	DropoffLocationLng    float64    `json:"dropoffLocationLng"`
	DropoffLocationLat    float64    `json:"dropoffLocationLat"`
	CreatedAt             time.Time  `json:"createdAt"`
	DeletedAt             *time.Time `json:"deletedAt,omitempty"`
}

func ListTripByRiderId(riderId int32) ([]*ListTripResp, error) {
	var trips []*ListTripResp
	rows, err := DBClient.pgPool.Query(context.Background(), listTripByRiderIdSQL, riderId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var trip ListTripResp
		if err := rows.Scan(
			&trip.Id,
			&trip.RiderId,
			&trip.DriverId,
			&trip.RequestId,
			&trip.RouteId,
			&trip.DriverName,
			&trip.DriverPictureUrl,
			&trip.DriverCarType,
			&trip.DriverCarPlate,
			&trip.RouteStartTime,
			&trip.PickupStartTime,
			&trip.PickupEndTime,
			&trip.RouteEndTime,
			&trip.RouteStartLocationLng,
			&trip.RouteStartLocationLat,
			&trip.RouteEndLocationLng,
			&trip.RouteEndLocationLat,
			&trip.PickupLocationLng,
			&trip.PickupLocationLat,
			&trip.DropoffLocationLng,
			&trip.DropoffLocationLat,
			&trip.CreatedAt,
			&trip.DeletedAt,
		); err != nil {
			return nil, ErrUndefined.WithCustomMessage(err.Error())
		}
		trips = append(trips, &trip)
	}
	return trips, nil
}

const listTripByDriverIdSQL = `
	SELECT 
		t.id,
		t.rider_id,
		t.driver_id,
		t.request_id,
		t.route_id,
		u.name,
		u.picture_url,
		u.car_type,
		u.car_plate,
		rout.start_time,
		req.pickup_start_time,
		req.pickup_end_time,
		rout.end_time,
		ST_X(rout.start_location), ST_Y(rout.start_location), 
		ST_X(rout.end_location), ST_Y(rout.end_location), 
		ST_X(req.pickup_location), ST_Y(req.pickup_location),
		ST_X(req.dropoff_location), ST_Y(req.dropoff_location),
		t.created_at,
		t.deleted_at
	FROM trips t
		JOIN users u ON t.driver_id = u.id
		JOIN requests req ON t.request_id = req.id
		JOIN routes rout ON t.route_id = rout.id
	WHERE t.rider_id = $1 AND t.deleted_at IS NULL;
`

func ListTripByDriverId(driverId int32) ([]*ListTripResp, error) {
	var trips []*ListTripResp
	rows, err := DBClient.pgPool.Query(context.Background(), listTripByDriverIdSQL, driverId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var trip ListTripResp
		if err := rows.Scan(
			&trip.Id,
			&trip.RiderId,
			&trip.DriverId,
			&trip.RequestId,
			&trip.RouteId,
			&trip.DriverName,
			&trip.DriverPictureUrl,
			&trip.DriverCarType,
			&trip.DriverCarPlate,
			&trip.RouteStartTime,
			&trip.PickupStartTime,
			&trip.PickupEndTime,
			&trip.RouteEndTime,
			&trip.RouteStartLocationLng,
			&trip.RouteStartLocationLat,
			&trip.RouteEndLocationLng,
			&trip.RouteEndLocationLat,
			&trip.PickupLocationLng,
			&trip.PickupLocationLat,
			&trip.DropoffLocationLng,
			&trip.DropoffLocationLat,
			&trip.CreatedAt,
			&trip.DeletedAt,
		); err != nil {
			return nil, ErrUndefined.WithCustomMessage(err.Error())
		}
		trips = append(trips, &trip)
	}
	return trips, nil
}

const getTripSQL = `
	SELECT *
	FROM trips
	WHERE id = $1 AND deleted_at IS NULL;
`

func GetTrip(id int32) (*model.Trip, error) {
	var trip model.Trip
	if err := DBClient.pgPool.QueryRow(context.Background(), getTripSQL, id).Scan(
		&trip.Id,
		&trip.RiderId,
		&trip.DriverId,
		&trip.RequestId,
		&trip.RouteId,
		&trip.CreatedAt,
		&trip.DeletedAt,
	); err != nil {
		return nil, Match(err, pgx.ErrNoRows, ErrTripNotFound).Return()
	}
	return &trip, nil
}

const createTripSQL = `
	INSERT INTO trips (rider_id, driver_id, request_id, route_id)
	VALUES (
		$1,
		$2, 
		$3, 
		$4
	)
	RETURNING id, created_at;
`

func CreateTrip(trip *model.Trip) (*model.Trip, error) {
	if err := DBClient.pgPool.QueryRow(context.Background(), createTripSQL,
		trip.RiderId,
		trip.DriverId,
		trip.RequestId,
		trip.RouteId,
	).Scan(
		&trip.Id,
		&trip.CreatedAt,
	); err != nil {
		return nil, ErrUndefined.WithCustomMessage(err.Error())
	}
	return trip, nil
}
