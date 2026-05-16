package handler

import (
	"auth-service/application/dto"
	"auth-service/application/usecase"
	pkgResponse "auth-service/pkg/response"
	"net/http"
)

type ListUserHandler struct {
	listUserUseCase *usecase.ListUserUseCase
}

func NewListUserHandler(listUserUseCase *usecase.ListUserUseCase) *ListUserHandler {
	return &ListUserHandler{listUserUseCase}
}

func (h *ListUserHandler) Execute(resp http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	reqQuery := dto.NewListUserFilterRequest()
	if err := reqQuery.Validate(req); err != nil {
		pkgResponse.ErrorWithCode(resp, http.StatusUnprocessableEntity, "failed to validate request query", err)
		return
	}

	data, err := h.listUserUseCase.Execute(ctx, reqQuery)
	if err != nil {
		pkgResponse.Error(resp, "failed to get user list", err)
		return
	}

	pkgResponse.SuccessWithMetaPagination(resp, "successfully get users data", data, reqQuery.PaginationRequest)
}
