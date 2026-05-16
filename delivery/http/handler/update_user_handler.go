package handler

import (
	"auth-service/application/usecase"
	"auth-service/delivery/http/middlewares"
	"auth-service/delivery/http/request"
	"auth-service/pkg/helpers"
	pkgResponse "auth-service/pkg/response"
	"net/http"

	"github.com/gorilla/mux"
)

type UpdateUserHandler struct {
	updateUserUseCase *usecase.UpdateUserUseCase
}

func NewUpdateUserHandler(updateUserUseCase *usecase.UpdateUserUseCase) *UpdateUserHandler {
	return &UpdateUserHandler{
		updateUserUseCase: updateUserUseCase,
	}
}

func (h *UpdateUserHandler) Execute(resp http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	idParam := mux.Vars(req)["id"]
	parsedID, err := helpers.ParseUintParam(idParam, "UUSUC", "001", "invalid mandatory parameter")
	if err != nil {
		pkgResponse.Error(resp, "failed to update user", err)
		return
	}
	body := ctx.Value(middlewares.BodyKey).(request.UpdateUserRequest)

	data, err := h.updateUserUseCase.Execute(ctx, parsedID, &body)
	if err != nil {
		pkgResponse.Error(resp, "failed to update user", err)
		return
	}

	pkgResponse.Success(resp, "successfully update user data", data)
}
