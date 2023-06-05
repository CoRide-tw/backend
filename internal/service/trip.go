package service

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/CoRide-tw/backend/internal/db"
	"github.com/CoRide-tw/backend/internal/model"
	"github.com/gin-gonic/gin"
)

type tripSvc struct {
}

func (s *tripSvc) List(c *gin.Context) {
	stringId, idExist := c.GetQuery("userId")
	if !idExist {
		c.JSON(http.StatusBadRequest, gin.H{"error": errors.New("missing id in query params")})
		return
	}
	userId, err := strconv.Atoi(stringId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	riderTrips, riderRespErr := db.ListTripByRiderId(int32(userId))
	if riderRespErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": riderRespErr.Error()})
		return
	}

	driverTrips, driverRespErr := db.ListTripByDriverId(int32(userId))
	if driverRespErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": driverRespErr.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"rider":  riderTrips,
		"driver": driverTrips,
	})
}

func (s *tripSvc) Get(c *gin.Context) {
	stringId := c.Param("id")
	tripId, err := strconv.Atoi(stringId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// get route from db
	trip, err := db.GetTrip(int32(tripId))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, trip)
}

func (s *tripSvc) Create(c *gin.Context) {
	var trip model.Trip
	if err := c.ShouldBindJSON(&trip); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// create route in db
	tripResp, err := db.CreateTrip(&trip)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, tripResp)
}
