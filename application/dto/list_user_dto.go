package dto

import (
	pkgPagination "auth-service/pkg/pagination"
	"net/http"
)

type ListUserFilterRequest struct {
	*pkgPagination.PaginationRequest
	IsAdmin *bool `json:"isAdmin"`
}

// NewListUserFilterRequest initialises the struct with safe pagination defaults.
func NewListUserFilterRequest() *ListUserFilterRequest {
	return &ListUserFilterRequest{
		PaginationRequest: pkgPagination.NewPaginationRequest(),
	}
}

func (f *ListUserFilterRequest) Validate(req *http.Request) error {
	if err := f.PaginationRequest.ParseFromQuery(req); err != nil {
		return err
	}

	if val := req.URL.Query().Get("isAdmin"); val != "" {
		b := val == "true"
		f.IsAdmin = &b
	}

	return nil
}
