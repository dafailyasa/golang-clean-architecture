package handler

import (
	"auth-service/internal/application/usecase"
	"auth-service/pkg/helpers"
	pkgResponse "auth-service/pkg/response"
	"net/http"

	//"github.com/gorilla/mux"

	"github.com/gorilla/mux"
)

type DetailGetHandler struct {
	detailUserUseCase *usecase.GetUserUseCase
}

func NewGetUserHandler(detailUserUseCase *usecase.GetUserUseCase) *DetailGetHandler {
	return &DetailGetHandler{detailUserUseCase}
}

func (h *DetailGetHandler) Execute(resp http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	idParam := mux.Vars(req)["id"]
	parsedID, err := helpers.ParseUintParam(idParam, "DUSUC", "001", "invalid mandatory parameter")
	if err != nil {
		pkgResponse.Error(resp, "failed to get detail user", err)
		return
	}

	data, err := h.detailUserUseCase.Execute(ctx, parsedID)
	if err != nil {
		pkgResponse.Error(resp, "failed to get detail user", err)
		return
	}

	pkgResponse.Success(resp, "successfully get detail user", data)
	return
}
