package handler

import (
	"auth-service/internal/application/dto"
	"auth-service/internal/application/usecase"
	"auth-service/internal/presentation/http/middlewares"
	pkgResponse "auth-service/pkg/response"
	"net/http"
)

type RefreshTokenUserHandler struct {
	refreshTokenUserUserUseCase *usecase.RefreshTokenUserUserUseCase
}

func NewRefreshTokenUserHandler(refreshTokenUserUserUseCase *usecase.RefreshTokenUserUserUseCase) *RefreshTokenUserHandler {
	return &RefreshTokenUserHandler{
		refreshTokenUserUserUseCase: refreshTokenUserUserUseCase,
	}
}

func (h *RefreshTokenUserHandler) Execute(resp http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	body := ctx.Value(middlewares.BodyKey).(dto.RefreshTokenUserDTO)

	data, err := h.refreshTokenUserUserUseCase.Execute(ctx, &body)
	if err != nil {
		pkgResponse.Error(resp, "failed generate new token from refresh token", err)
		return
	}

	pkgResponse.Success(resp, "successfully generate new token from refresh token", data)
}
