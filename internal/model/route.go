package model

import "time"

type Route struct {
	Id        int32      `json:"id"`
	DriverId  int32      `json:"driverId"`
	StartLong float64    `json:"startLong"`
	StartLat  float64    `json:"startLat"`
	EndLong   float64    `json:"endLong"`
	EndLat    float64    `json:"endLat"`
	StartTime time.Time  `json:"startTime"`
	EndTime   time.Time  `json:"endTime"`
	Capacity  int32      `json:"capacity"`
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt time.Time  `json:"updatedAt"`
	DeletedAt *time.Time `json:"deletedAt,omitempty"`
}
