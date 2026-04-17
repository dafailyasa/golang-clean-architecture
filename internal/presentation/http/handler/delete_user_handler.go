package handler

import (
	"auth-service/internal/application/usecase"
	"auth-service/pkg/helpers"
	pkgResponse "auth-service/pkg/response"
	"net/http"

	"github.com/gorilla/mux"
)

type DeleteUserHandler struct {
	deleteUserUseCase *usecase.DeleteUserUseCase
}

func NewDeleteUserHandler(deleteUserUseCase *usecase.DeleteUserUseCase) *DeleteUserHandler {
	return &DeleteUserHandler{deleteUserUseCase: deleteUserUseCase}
}

func (h *DeleteUserHandler) Execute(resp http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	idParam := mux.Vars(req)["id"]
	parsedID, err := helpers.ParseUintParam(idParam, "DUSUC", "001", "invalid mandatory parameter")
	if err != nil {
		pkgResponse.Error(resp, "failed to delete user", err)
		return
	}

	err = h.deleteUserUseCase.Execute(ctx, parsedID)
	if err != nil {
		pkgResponse.Error(resp, "failed to delete user", err)
		return
	}

	pkgResponse.SuccessWithNoContent(resp, "successfully delete user data")
	return
}
