package service

type Service struct {
	User  *userSvc
	Route *routeSvc
}

func NewService() *Service {
	return &Service{
		User:  &userSvc{},
		Route: &routeSvc{},
	}
}
