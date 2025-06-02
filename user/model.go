package user

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func Equal(u1, u2 User) bool {
	return u1.Username == u2.Username && u1.Password == u2.Password
}
