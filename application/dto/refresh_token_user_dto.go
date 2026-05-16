package dto

type RefreshTokenUserDTO struct {
	RefreshToken string `json:"refreshToken" validate:"required,jwtToken"`
}
