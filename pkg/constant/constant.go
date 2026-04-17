package constant

const (
	KeycloakAdminCLIConst              = "admin-cli"
	KeycloakGrantTypePasswordConst     = "password"
	KeycloakGrantTypeRefreshTokenConst = "refresh_token"
	KeycloakScope                      = "openid profile email"
	KeycloakClientCredentialsConst     = "client_credentials"
	ApplicationJsonConst               = "application/json"
	FormUrlEncodedConst                = "application/x-www-form-urlencoded"
)

// pagination
const (
	DefaultPage      = 1
	DefaultLimit     = 10
	DefaultOrder     = "created_at"
	DefaultSortOrder = "desc"
)
