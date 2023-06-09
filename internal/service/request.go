package service

import (
	"go.uber.org/zap"
	"net/http"
	"strconv"

	"github.com/CoRide-tw/backend/internal/constants"
	"github.com/CoRide-tw/backend/internal/db"
	"github.com/CoRide-tw/backend/internal/model"
	"github.com/CoRide-tw/backend/internal/util"
	"github.com/gin-gonic/gin"
)

type requestSvc struct {
	Logger *zap.SugaredLogger
}

func (s *requestSvc) List(c *gin.Context) {
	parsedQuery, err := util.ParseListRequestQuery(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if parsedQuery.RiderId != 0 {
		requests, err := db.ListRequestsByRiderId(parsedQuery.RiderId)
		if err != nil {
			s.Logger.Error(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, requests)
		return
	}

	if parsedQuery.RouteId != 0 {
		requests, err := db.ListRequestsByRouteId(parsedQuery.RouteId)
		if err != nil {
			s.Logger.Error(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, requests)
		return
	}
}

func (s *requestSvc) Get(c *gin.Context) {
	stringId := c.Param("id")
	requestId, err := strconv.Atoi(stringId)
	if err != nil {
		s.Logger.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// get route from db
	request, err := db.GetRequest(int32(requestId))
	if err != nil {
		s.Logger.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, request)
}

func (s *requestSvc) Create(c *gin.Context) {
	var request model.Request
	if err := c.ShouldBindJSON(&request); err != nil {
		s.Logger.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// create route in db
	requestResp, err := db.CreateRequest(&request)
	if err != nil {
		s.Logger.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, requestResp)
}

func (s *requestSvc) Deny(c *gin.Context) {
	stringId := c.Param("id")
	requestId, err := strconv.Atoi(stringId)
	if err != nil {
		s.Logger.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := db.UpdateRequestStatus(int32(requestId), constants.RequestStatusDenied); err != nil {
		s.Logger.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{})
}

func (s *requestSvc) Delete(c *gin.Context) {
	stringId := c.Param("id")
	requestId, err := strconv.Atoi(stringId)
	if err != nil {
		s.Logger.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := db.DeleteRequest(int32(requestId)); err != nil {
		s.Logger.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{})
}
