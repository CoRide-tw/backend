package service

type Service struct {
	User *userSvc
}

func NewService() *Service {
	return &Service{
		User: &userSvc{},
	}
}
