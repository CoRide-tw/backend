package service

import (
	"context"
	"net/http"

	"github.com/CoRide-tw/backend/internal/config"
	"github.com/gin-gonic/gin"
	"googlemaps.github.io/maps"
)

type googleApiSvc struct {
}

// search latlng with text
func (s *googleApiSvc) GetGeocodingWithTextSearch(c *gin.Context) {
	text, queryExist := c.GetQuery("text")
	if !queryExist {
		c.JSON(http.StatusBadRequest, gin.H{"error": "text is required"})
		return
	}

	mapsClient, err := maps.NewClient(maps.WithAPIKey(config.Env.GoogleMapsApiKey))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	res, err := mapsClient.FindPlaceFromText(context.Background(), &maps.FindPlaceFromTextRequest{
		Input:     text,
		Language:  "zh-TW",
		Fields:    []maps.PlaceSearchFieldMask{maps.PlaceSearchFieldMaskGeometry, maps.PlaceSearchFieldMaskFormattedAddress},
		InputType: maps.FindPlaceFromTextInputTypeTextQuery,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if len(res.Candidates) == 0 {
		c.JSON(http.StatusOK, gin.H{"lng": "", "lat": "", "address": ""})
		return
	}

	c.JSON(http.StatusOK, gin.H{"lng": res.Candidates[0].Geometry.Location.Lng, "lat": res.Candidates[0].Geometry.Location.Lat, "address": res.Candidates[0].FormattedAddress})
}

// query placeId, description with text
func (s *googleApiSvc) GetPlaceAutocomplete(c *gin.Context) {
	place, queryExist := c.GetQuery("place")
	if !queryExist {
		c.JSON(http.StatusBadRequest, gin.H{"error": "place is required"})
		return
	}

	mapsClient, err := maps.NewClient(maps.WithAPIKey(config.Env.GoogleMapsApiKey))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	res, err := mapsClient.PlaceAutocomplete(context.Background(), &maps.PlaceAutocompleteRequest{
		Input:    place,
		Language: "zh-TW",
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if len(res.Predictions) == 0 {
		c.JSON(http.StatusOK, gin.H{"predictions": "", "placeId": ""})
		return
	}

	c.JSON(http.StatusOK, gin.H{"predictions": res.Predictions[0].Description, "placeId": res.Predictions[0].PlaceID})
}
