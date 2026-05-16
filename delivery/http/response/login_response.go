package response

type UserLoginData struct {
	FullName string `json:"fullName"`
}

type LoginUserResponse struct {
	AccessToken  string        `json:"accessToken"`
	RefreshToken string        `json:"refreshToken"`
	User         UserLoginData `json:"user"`
}

func NewLoginUserResponse(accessToken, refreshToken, name string) *LoginUserResponse {
	return &LoginUserResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User: UserLoginData{
			FullName: name,
		},
	}
}
