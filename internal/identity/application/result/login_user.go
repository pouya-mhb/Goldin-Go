package result

// LoginUser is returned after credentials are authenticated.
type LoginUser struct {
	UserID string
	Email  string
}
