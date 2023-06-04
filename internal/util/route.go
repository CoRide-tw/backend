package util

import (
	"log"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type ParsedListNearestRoutesQuery struct {
	StartLong       float64
	StartLat        float64
	EndLong         float64
	EndLat          float64
	PickupStartTime time.Time
	PickupEndTime   time.Time
}

func ParseListNearestRoutesQuery(c *gin.Context) (*ParsedListNearestRoutesQuery, error) {
	var parsedQuery ParsedListNearestRoutesQuery
	var err error

	parsedQuery.StartLong, err = strconv.ParseFloat(c.Query("startLong"), 64)
	if err != nil {
		return nil, err
	}
	parsedQuery.StartLat, err = strconv.ParseFloat(c.Query("startLat"), 64)
	if err != nil {
		return nil, err
	}
	parsedQuery.EndLong, err = strconv.ParseFloat(c.Query("endLong"), 64)
	if err != nil {
		return nil, err
	}
	parsedQuery.EndLat, err = strconv.ParseFloat(c.Query("endLat"), 64)
	if err != nil {
		return nil, err
	}
	parsedQuery.PickupStartTime, err = time.Parse(time.RFC3339, c.Query("startTime"))
	if err != nil {
		log.Println("PickupStartTime", err)
		return nil, err
	}
	parsedQuery.PickupEndTime, err = time.Parse(time.RFC3339, c.Query("endTime"))
	if err != nil {
		return nil, err
	}
	return &parsedQuery, nil
}
