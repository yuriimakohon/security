package user

type Repository interface {
	Create(user User) error
	Get(login string) (User, error)
	UpdateUsername(login, username string) error
}
