package keycloack

type RegisterUserDTO struct {
	Username      string           `json:"username"`
	Enabled       bool             `json:"enabled"`
	EmailVerified bool             `json:"emailVerified"`
	FirstName     string           `json:"firstName"`
	LastName      string           `json:"lastName"`
	Email         string           `json:"email"`
	Credentials   []CredentialsDTO `json:"credentials"`
}

type CredentialsDTO struct {
	Type      string `json:"type"`
	Value     string `json:"value"`
	Temporary bool   `json:"temporary"`
}

type RegisterUserResponseDTO struct {
	BaseErrorKeycloakDTO
}

type AccessTokenDTO struct {
	ClientID     string `json:"client_id"`
	Username     string `json:"username"`
	Password     string `json:"password"`
	GrantType    string `json:"grant_type"`
	ClientSecret string `json:"client_secret"`
	Scope        string `json:"scope,omitempty"`
	//RefreshToken string `json:"refresh_token,omitempty"`
}

type RefreshTokenDTO struct {
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	RefreshToken string `json:"refresh_token"`
	GrantType    string `json:"grant_type"`
}

type AccessTokenResponseDTO struct {
	AccessToken           string `json:"access_token"`
	ExpiresIn             int    `json:"expires_in"`
	RefreshToken          string `json:"refresh_token"`
	RefreshTokenExpiredIn int    `json:"refresh_expires_in"`
	TokenType             string `json:"token_type"`
	Error                 string `json:"error"`
	ErrorDescription      string `json:"error_description"`
}

type BaseErrorKeycloakDTO struct {
	Error            string `json:"error,omitempty"`
	ErrorDescription string `json:"error_description,omitempty"`
	ErrorMessage     string `json:"errorMessage,omitempty"`
}

type UserInfoResponseDTO struct {
	Sub               string `json:"sub"`
	Email             string `json:"email"`
	EmailVerified     bool   `json:"email_verified"`
	PreferredUsername string `json:"preferred_username"`
	GivenName         string `json:"given_name"`
	FamilyName        string `json:"family_name"`
}

type UpdateUserDTO struct {
	FirstName     string           `json:"firstName,omitempty"`
	LastName      string           `json:"lastName,omitempty"`
	Email         string           `json:"email,omitempty"`
	Enabled       bool             `json:"enabled"`
	EmailVerified bool             `json:"emailVerified"`
	Credentials   []CredentialsDTO `json:"credentials,omitempty"`
}
