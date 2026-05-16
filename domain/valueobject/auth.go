package valueobject

type RegisterUserParam struct {
	Username  string
	FirstName string
	LastName  string
	Email     string
	Password  string
}

type UpdateUserParam struct {
	FirstName string
	LastName  string
	Email     string
	Password  string
	Enabled   bool
}

type TokenResponse struct {
	AccessToken  string
	RefreshToken string
}

type UserInfoResponse struct {
	Sub   string
	Email string
}
