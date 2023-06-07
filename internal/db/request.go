package db

import (
	"context"
	"time"

	"github.com/CoRide-tw/backend/internal/constants"
	. "github.com/CoRide-tw/backend/internal/errors/generated/dberr"
	"github.com/CoRide-tw/backend/internal/model"
	. "github.com/DenChenn/blunder/pkg/blunder"
	"github.com/jackc/pgx/v5"
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
		tips INT NOT NULL,
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
		tips,
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
		&request.Tips,
		&request.Status,
		&request.CreatedAt,
		&request.UpdatedAt,
	); err != nil {
		return nil, Match(err, pgx.ErrNoRows, ErrRequestNotFound).Return()
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
		tips,
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
			&request.Tips,
			&request.Status,
			&request.CreatedAt,
			&request.UpdatedAt,
		); err != nil {
			return nil, ErrUndefined.WithCustomMessage(err.Error())
		}
		requests = append(requests, &request)
	}
	return requests, nil
}

const listRequestsByRouteIdSQL = `
	SELECT 
		r.id,
		r.rider_id,
		r.route_id,
		ST_X(pickup_location), ST_Y(pickup_location),
		ST_X(dropoff_location), ST_Y(dropoff_location),
		r.pickup_start_time, 
		r.pickup_end_time,
		r.tips,
		r.status,
		u.name,
		u.picture_url,
		r.created_at, 
		r.updated_at
	FROM requests r
		JOIN users u ON r.rider_id = u.id
	WHERE r.route_id = $1 AND r.deleted_at IS NULL;
`

type ListRequestsByRouteIdResp struct {
	Id              int32      `json:"id"`
	RiderId         int32      `json:"riderId"`
	RouteId         int32      `json:"routeId"`
	PickupLong      float64    `json:"pickupLong"`
	PickupLat       float64    `json:"pickupLat"`
	DropoffLong     float64    `json:"dropoffLong"`
	DropoffLat      float64    `json:"dropoffLat"`
	PickupStartTime time.Time  `json:"pickupStartTime"`
	PickupEndTime   time.Time  `json:"pickupEndTime"`
	Tips            int32      `json:"tips"`
	Status          string     `json:"status"`
	RiderName       string     `json:"riderName"`
	RiderPictureUrl string     `json:"riderPictureUrl"`
	CreatedAt       time.Time  `json:"createdAt"`
	UpdatedAt       time.Time  `json:"updatedAt"`
	DeletedAt       *time.Time `json:"deletedAt,omitempty"`
}

func ListRequestsByRouteId(routeId int32) ([]*ListRequestsByRouteIdResp, error) {
	rows, err := DBClient.pgPool.Query(context.Background(), listRequestsByRouteIdSQL, routeId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var requests []*ListRequestsByRouteIdResp
	for rows.Next() {
		var request ListRequestsByRouteIdResp
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
			&request.Tips,
			&request.Status,
			&request.RiderName,
			&request.RiderPictureUrl,
			&request.CreatedAt,
			&request.UpdatedAt,
		); err != nil {
			return nil, ErrUndefined.WithCustomMessage(err.Error())
		}
		requests = append(requests, &request)
	}
	return requests, nil
}

const createRequestSQL = `
	INSERT INTO requests (rider_id, route_id, pickup_location, dropoff_location, pickup_start_time, pickup_end_time, tips, status)
	VALUES (
		$1,
		$2,
		ST_SetSRID(ST_MakePoint($3, $4), 4326),
		ST_SetSRID(ST_MakePoint($5, $6), 4326),
		$7,
		$8,
		$9,
		$10
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
		request.Tips,
		constants.RequestStatusPending,
	).Scan(
		&request.Id,
		&request.Status,
		&request.CreatedAt,
		&request.UpdatedAt,
	); err != nil {
		return nil, ErrUndefined.WithCustomMessage(err.Error())
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
		return ErrUndefined.WithCustomMessage(err.Error())
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
		return ErrUndefined.WithCustomMessage(err.Error())
	}
	return nil
}
