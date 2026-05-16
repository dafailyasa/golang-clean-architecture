package handler

import (
	"auth-service/application/usecase"
	"auth-service/delivery/http/middlewares"
	"auth-service/delivery/http/request"
	pkgResponse "auth-service/pkg/response"
	"net/http"
)

type CreateUserHandler struct {
	createUserUseCase *usecase.CreateUserUseCase
}

func NewCreateUserHandler(createUserUseCase *usecase.CreateUserUseCase) *CreateUserHandler {
	return &CreateUserHandler{
		createUserUseCase: createUserUseCase,
	}
}

func (h *CreateUserHandler) Execute(resp http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	body := ctx.Value(middlewares.BodyKey).(request.CreateUserRequest)

	data, err := h.createUserUseCase.Execute(ctx, &body)
	if err != nil {
		pkgResponse.Error(resp, "failed to create user", err)
		return
	}

	pkgResponse.Created(resp, "successfully create user data", data)
	return
}
