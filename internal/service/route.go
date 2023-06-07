package service

import (
	"go.uber.org/zap"
	"net/http"
	"strconv"

	"github.com/CoRide-tw/backend/internal/db"
	"github.com/CoRide-tw/backend/internal/model"
	"github.com/CoRide-tw/backend/internal/util"
	"github.com/gin-gonic/gin"
)

type routeSvc struct {
	Logger *zap.SugaredLogger
}

func (s *routeSvc) ListNearestRoutes(c *gin.Context) {
	parsedQuery, err := util.ParseListNearestRoutesQuery(c)
	if err != nil {
		s.Logger.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// get nearest routes from db
	routes, err := db.ListNearestRoutes(parsedQuery.StartLong, parsedQuery.StartLat, parsedQuery.EndLong, parsedQuery.EndLat, parsedQuery.PickupStartTime, parsedQuery.PickupEndTime)
	if err != nil {
		s.Logger.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, routes)
}

func (s *routeSvc) Get(c *gin.Context) {
	stringId := c.Param("id")
	routeId, err := strconv.Atoi(stringId)
	if err != nil {
		s.Logger.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// get route from db
	route, err := db.GetRoute(int32(routeId))
	if err != nil {
		s.Logger.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, route)
}

func (s *routeSvc) Create(c *gin.Context) {
	var route model.Route
	if err := c.ShouldBindJSON(&route); err != nil {
		s.Logger.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// create route in db
	routeResp, err := db.CreateRoute(&route)
	if err != nil {
		s.Logger.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, routeResp)
}

func (s *routeSvc) Delete(c *gin.Context) {
	stringId := c.Param("id")
	routeId, err := strconv.Atoi(stringId)
	if err != nil {
		s.Logger.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := db.DeleteRoute(int32(routeId)); err != nil {
		s.Logger.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{})
}
