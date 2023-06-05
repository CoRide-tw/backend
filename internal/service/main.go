package service

type Service struct {
	User    *userSvc
	Route   *routeSvc
	Request *requestSvc
	Trip    *tripSvc
}

func NewService() *Service {
	return &Service{
		User:    &userSvc{},
		Route:   &routeSvc{},
		Request: &requestSvc{},
		Trip:    &tripSvc{},
	}
}
