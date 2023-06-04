package service

import (
	"net/http"
	"strconv"

	"github.com/CoRide-tw/backend/internal/db"
	"github.com/CoRide-tw/backend/internal/model"
	"github.com/CoRide-tw/backend/internal/util"
	"github.com/gin-gonic/gin"
)

type routeSvc struct {
}

func (s *routeSvc) ListNearestRoutes(c *gin.Context) {
	parsedQuery, err := util.ParseListNearestRoutesQuery(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// get nearest routes from db
	routes, err := db.ListNearestRoutes(parsedQuery.StartLong, parsedQuery.StartLat, parsedQuery.EndLong, parsedQuery.EndLat, parsedQuery.PickupStartTime, parsedQuery.PickupEndTime)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, routes)
}

func (s *routeSvc) Get(c *gin.Context) {
	stringId := c.Param("id")
	routeId, err := strconv.Atoi(stringId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// get route from db
	route, err := db.GetRoute(int32(routeId))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, route)
}

func (s *routeSvc) Create(c *gin.Context) {
	var route model.Route
	if err := c.ShouldBindJSON(&route); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// create route in db
	routeResp, err := db.CreateRoute(&route)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, routeResp)
}

func (s *routeSvc) Update(c *gin.Context) {
	stringId := c.Param("id")
	routeId, err := strconv.Atoi(stringId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var route model.Route
	if err := c.ShouldBindJSON(&route); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if routeId != int(route.Id) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "route id in path and body are not the same"})
		return
	}

	// create route in db
	routeResp, err := db.UpdateRoute(&route)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, routeResp)
}

func (s *routeSvc) Delete(c *gin.Context) {
	stringId := c.Param("id")
	routeId, err := strconv.Atoi(stringId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := db.DeleteRoute(int32(routeId)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{})
}
