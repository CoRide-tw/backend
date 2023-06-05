package util

import (
	"errors"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ParsedListRequestQuery struct {
	RiderId int32
	RouteId int32
}

func ParseListRequestQuery(c *gin.Context) (*ParsedListRequestQuery, error) {
	var parsedQuery ParsedListRequestQuery

	stringRiderId, riderIdExist := c.GetQuery("riderId")
	stringRouteId, routeIdExist := c.GetQuery("routeId")

	if riderIdExist && routeIdExist {
		return nil, errors.New("cannot query by both riderId and routeId")
	}

	if riderIdExist {
		rawRiderId, riderIdErr := strconv.ParseInt(stringRiderId, 10, 32)
		if riderIdErr != nil {
			return nil, errors.New("invalid riderId")
		}
		parsedQuery.RiderId = int32(rawRiderId)
		return &parsedQuery, nil
	}

	if routeIdExist {
		rawRouteId, routeIdErr := strconv.ParseInt(stringRouteId, 10, 32)
		if routeIdErr != nil {
			return nil, errors.New("invalid routeId")
		}
		parsedQuery.RouteId = int32(rawRouteId)
		return &parsedQuery, nil
	}

	return nil, errors.New("must query by either riderId or routeId")
}
