package db

import (
	"context"

	"github.com/CoRide-tw/backend/internal/constants"
	"github.com/CoRide-tw/backend/internal/model"
)

const createRequestTable = `
	CREATE EXTENSION IF NOT EXISTS postgis;

	CREATE TABLE IF NOT EXISTS requests (
		id SERIAL,
		rider_id INT NOT NULL,
		route_id INT NOT NULL,
		pickup_location GEOMETRY(Point, 4326) NOT NULL,
		dropoff_location GEOMETRY(Point, 4326) NOT NULL,
		pickup_start_time TIMESTAMP WITH TIME ZONE NOT NULL,
		pickup_end_time TIMESTAMP WITH TIME ZONE NOT NULL,
		status VARCHAR(50) NOT NULL,
		created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL,
		updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL,
		deleted_at TIMESTAMP WITH TIME ZONE
	);

	CREATE INDEX IF NOT EXISTS requests_pickup_location_dropoff_location_idx 
		ON requests USING gist(pickup_location, dropoff_location);
`

func initRequestTable() error {
	if _, err := DBClient.pgPool.Exec(context.Background(), createRequestTable); err != nil {
		return err
	}
	return nil
}

const getRequestSQL = `
	SELECT
		id,
		rider_id,
		route_id,
		ST_X(pickup_location), ST_Y(pickup_location),
		ST_X(dropoff_location), ST_Y(dropoff_location),
		pickup_start_time, pickup_end_time,
		status,
		created_at, updated_at
	FROM requests
	WHERE id = $1 AND deleted_at IS NULL;
`

func GetRequest(id int32) (*model.Request, error) {
	var request model.Request
	if err := DBClient.pgPool.QueryRow(context.Background(), getRequestSQL, id).Scan(
		&request.Id,
		&request.RiderId,
		&request.RouteId,
		&request.PickupLong,
		&request.PickupLat,
		&request.DropoffLong,
		&request.DropoffLat,
		&request.PickupStartTime,
		&request.PickupEndTime,
		&request.Status,
		&request.CreatedAt,
		&request.UpdatedAt,
	); err != nil {
		return nil, err
	}
	return &request, nil
}

const listRequestsByRiderIdSQL = `
	SELECT 
		id,
		rider_id,
		route_id,
		ST_X(pickup_location), ST_Y(pickup_location),
		ST_X(dropoff_location), ST_Y(dropoff_location),
		pickup_start_time, 
		pickup_end_time,
		status,
		created_at, 
		updated_at
	FROM requests
	WHERE rider_id = $1 AND deleted_at IS NULL;
`

func ListRequestsByRiderId(riderId int32) ([]*model.Request, error) {
	rows, err := DBClient.pgPool.Query(context.Background(), listRequestsByRiderIdSQL, riderId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var requests []*model.Request
	for rows.Next() {
		var request model.Request
		if err := rows.Scan(
			&request.Id,
			&request.RiderId,
			&request.RouteId,
			&request.PickupLong,
			&request.PickupLat,
			&request.DropoffLong,
			&request.DropoffLat,
			&request.PickupStartTime,
			&request.PickupEndTime,
			&request.Status,
			&request.CreatedAt,
			&request.UpdatedAt,
		); err != nil {
			return nil, err
		}
		requests = append(requests, &request)
	}
	return requests, nil
}

const listRequestsByRouteIdSQL = `
	SELECT 
		id,
		rider_id,
		route_id,
		ST_X(pickup_location), ST_Y(pickup_location),
		ST_X(dropoff_location), ST_Y(dropoff_location),
		pickup_start_time, 
		pickup_end_time,
		status,
		created_at, 
		updated_at
	FROM requests
	WHERE route_id = $1 AND deleted_at IS NULL;
`

func ListRequestsByRouteId(routeId int32) ([]*model.Request, error) {
	rows, err := DBClient.pgPool.Query(context.Background(), listRequestsByRouteIdSQL, routeId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var requests []*model.Request
	for rows.Next() {
		var request model.Request
		if err := rows.Scan(
			&request.Id,
			&request.RiderId,
			&request.RouteId,
			&request.PickupLong,
			&request.PickupLat,
			&request.DropoffLong,
			&request.DropoffLat,
			&request.PickupStartTime,
			&request.PickupEndTime,
			&request.Status,
			&request.CreatedAt,
			&request.UpdatedAt,
		); err != nil {
			return nil, err
		}
		requests = append(requests, &request)
	}
	return requests, nil
}

const createRequestSQL = `
	INSERT INTO requests (rider_id, route_id, pickup_location, dropoff_location, pickup_start_time, pickup_end_time, status)
	VALUES (
		$1,
		$2,
		ST_SetSRID(ST_MakePoint($3, $4), 4326),
		ST_SetSRID(ST_MakePoint($5, $6), 4326),
		$7,
		$8,
		$9
	)
	RETURNING id, status, created_at, updated_at;
`

func CreateRequest(request *model.Request) (*model.Request, error) {
	if err := DBClient.pgPool.QueryRow(context.Background(), createRequestSQL,
		request.RiderId,
		request.RouteId,
		request.PickupLong,
		request.PickupLat,
		request.DropoffLong,
		request.DropoffLat,
		request.PickupStartTime,
		request.PickupEndTime,
		constants.RequestStatusPending,
	).Scan(
		&request.Id,
		&request.Status,
		&request.CreatedAt,
		&request.UpdatedAt,
	); err != nil {
		return nil, err
	}
	return request, nil
}

const updateRequestSQL = `
	UPDATE requests SET
		pickup_location = ST_SetSRID(ST_MakePoint($4, $5), 4326), 
		dropoff_location = ST_SetSRID(ST_MakePoint($6, $7), 4326), 
		pickup_start_time = $7,
		pickup_end_time = $8,
		status = $9,
		updated_at = NOW()
	WHERE id = $1, AND rider_id = $2 AND route_id = $3 AND deleted_at IS NULL
	RETURNING
		id,
		status,
		created_at, 
		updated_at;
`

func UpdateRequest(request *model.Request) (*model.Request, error) {
	if err := DBClient.pgPool.QueryRow(context.Background(), updateRequestSQL,
		&request.Id,
		&request.RiderId,
		&request.RouteId,
		&request.PickupLong,
		&request.PickupLat,
		&request.DropoffLong,
		&request.DropoffLat,
		&request.PickupStartTime,
		&request.PickupEndTime,
		constants.RequestStatusPending,
	).Scan(
		&request.Id,
		&request.Status,
		&request.CreatedAt,
		&request.UpdatedAt,
	); err != nil {
		return nil, err
	}
	return request, nil
}

const updateRequestStatusSQL = `
	UPDATE requests SET
		status = $2,
		updated_at = NOW()
	WHERE id = $1 AND deleted_at IS NULL
`

func UpdateRequestStatus(id int32, status string) error {
	if _, err := DBClient.pgPool.Exec(context.Background(), updateRequestStatusSQL,
		id, status,
	); err != nil {
		return err
	}
	return nil
}

const deleteRequestSQL = `
	UPDATE requests SET
		status = $2,
		deleted_at = NOW()
	WHERE id = $1 AND deleted_at IS NULL
`

func DeleteRequest(id int32) error {
	if _, err := DBClient.pgPool.Exec(context.Background(), deleteRequestSQL,
		id,
		constants.RequestStatusCancelled,
	); err != nil {
		return err
	}
	return nil
}
