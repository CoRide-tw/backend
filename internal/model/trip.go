package model

import "time"

type Trip struct {
	Id        int32      `json:"id"`
	RiderId   int32      `json:"riderId"`
	DriverId  int32      `json:"driverId"`
	RequestId int32      `json:"requestId"`
	RouteId   int32      `json:"routeId"`
	CreatedAt time.Time  `json:"createdAt"`
	DeletedAt *time.Time `json:"deletedAt,omitempty"`
}
