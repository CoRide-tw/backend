package model

import "time"

type Request struct {
	Id              int32      `json:"id"`
	RiderId         int32      `json:"riderId"`
	RouteId         int32      `json:"routeId"`
	PickupLong      float64    `json:"pickupLong"`
	PickupLat       float64    `json:"pickupLat"`
	DropoffLong     float64    `json:"dropoffLong"`
	DropoffLat      float64    `json:"dropoffLat"`
	PickupStartTime time.Time  `json:"pickupStartTime"`
	PickupEndTime   time.Time  `json:"pickupEndTime"`
	Status          string     `json:"status"`
	CreatedAt       time.Time  `json:"createdAt"`
	UpdatedAt       time.Time  `json:"updatedAt"`
	DeletedAt       *time.Time `json:"deletedAt,omitempty"`
}
