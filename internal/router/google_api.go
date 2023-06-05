package router

func (r *router) setGoogleApiRoutes() {
	googleApiRouter := r.Engine.Group("/google_api")

	googleApiRouter.GET("/geocoding", r.Service.GoogleApi.GetGeocodingWithTextSearch)
	googleApiRouter.GET("/place_autocomplete", r.Service.GoogleApi.GetPlaceAutocomplete)
}
