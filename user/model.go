package user

type User struct {
	Id       int64
	Username string
	Password string
}

func Equal(u1, u2 User) bool {
	return u1.Username == u2.Username && u1.Password == u2.Password
}
