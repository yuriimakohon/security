package user

type User struct {
	Username string `pg:"username"`
	Login    string `pg:"login"`
	Password string `pg:"password"`
}
