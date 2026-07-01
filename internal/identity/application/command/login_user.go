package command

// LoginUser carries credentials required to authenticate a user.
type LoginUser struct {
	Email    string
	Password string
}
