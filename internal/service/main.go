package service

import "go.uber.org/zap"

type Service struct {
	User      *userSvc
	Route     *routeSvc
	Request   *requestSvc
	Trip      *tripSvc
	GoogleApi *googleApiSvc
	Logger    *zap.SugaredLogger
}

func NewService(logger *zap.SugaredLogger) *Service {
	return &Service{
		User:    &userSvc{Logger: logger},
		Route:   &routeSvc{Logger: logger},
		Request: &requestSvc{Logger: logger},
		Trip:    &tripSvc{Logger: logger},
	}
}
