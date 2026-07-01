package result

// LoginUser is returned after credentials are authenticated.
type LoginUser struct {
	UserID                string
	Email                 string
	AccessToken           string
	RefreshToken          string
	TokenType             string
	AccessTokenExpiresIn  int64
	RefreshTokenExpiresIn int64
}
