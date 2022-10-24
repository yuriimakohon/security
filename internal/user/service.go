package user

import "golang.org/x/crypto/bcrypt"

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

// Create a new user if not exists. Password will be encrypted and salted with bcrypt algorithm
func (s *Service) Create(user User) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user.Password = string(hashedPassword)
	return s.repo.Create(user)
}

func (s *Service) Get(login string) (User, error) {
	return s.repo.Get(login)
}

// CheckPassword checks if user exists and password is correct
func (s *Service) CheckPassword(login, password string) (bool, error) {
	user, err := s.repo.Get(login)
	if err != nil {
		return false, err
	}

	return bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)) == nil, nil
}

func (s *Service) UpdateUsername(login, username string) error {
	return s.repo.UpdateUsername(login, username)
}
