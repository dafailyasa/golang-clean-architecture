package request

type RefreshTokenRequest struct {
	RefreshToken string `json:"refreshToken" validate:"required,jwtToken"`
}
