package handler

import (
	"auth-service/application/usecase"
	"auth-service/delivery/http/middlewares"
	"auth-service/delivery/http/request"
	pkgResponse "auth-service/pkg/response"
	"net/http"
)

type LoginUserHandler struct {
	loginUserUseCase *usecase.LoginUserUseCase
}

func NewLoginUserHandler(loginUserUseCase *usecase.LoginUserUseCase) *LoginUserHandler {
	return &LoginUserHandler{
		loginUserUseCase: loginUserUseCase,
	}
}

func (h *LoginUserHandler) Execute(resp http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	body := ctx.Value(middlewares.BodyKey).(request.LoginUserRequest)

	data, err := h.loginUserUseCase.Execute(ctx, &body)
	if err != nil {
		pkgResponse.Error(resp, "failed to login user", err)
		return
	}

	pkgResponse.Success(resp, "successfully login user", data)
	return
}
