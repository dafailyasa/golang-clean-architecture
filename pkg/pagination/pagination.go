package pagination

import (
	pkgConst "auth-service/pkg/constant"
	"errors"
	"math"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-playground/validator/v10"
)

var Validate = validator.New()

type PaginationRequest struct {
	Page      int    `json:"page"`
	Limit     int    `json:"limit"`
	SortBy    string `json:"sortBy,omitempty"`
	SortOrder string `json:"sortOrder,omitempty"`
	Search    string `json:"search,omitempty"`
	Total     int64  `json:"total"`
}

// NewPaginationRequest returns a *PaginationRequest with safe defaults applied.
func NewPaginationRequest() *PaginationRequest {
	return &PaginationRequest{
		Page:      pkgConst.DefaultPage,
		Limit:     pkgConst.DefaultLimit,
		SortBy:    pkgConst.DefaultOrder,
		SortOrder: pkgConst.DefaultSortOrder,
	}
}

// ParseFromQuery reads pagination parameters from the URL query string and
// overrides the current values. Unknown keys are silently ignored so that
// embed structs can handle their own extra fields.
func (p *PaginationRequest) ParseFromQuery(req *http.Request) error {
	q := req.URL.Query()

	if page := q.Get("page"); page != "" {
		val, err := strconv.Atoi(page)
		if err != nil || val < 1 {
			return errors.New("invalid page must be a positive integer")
		}
		p.Page = val
	}

	if limit := q.Get("limit"); limit != "" {
		val, err := strconv.Atoi(limit)
		if err != nil || val < 1 {
			return errors.New("invalid limit must be a positive integer")
		}
		p.Limit = val
	}

	if sort := q.Get("sortBy"); sort != "" {
		p.SortBy = sort
	}

	if order := q.Get("sortOrder"); order != "" {
		o := strings.ToLower(order)
		if o != "asc" && o != "desc" {
			return errors.New("invalid sortOrder must be 'asc' or 'desc'")
		}
		p.SortOrder = o
	}

	if search := q.Get("search"); search != "" {
		p.Search = search
	}

	return nil
}

// ---------------------------------------------------------------------------
// Accessor helpers
// ---------------------------------------------------------------------------

func (p *PaginationRequest) GetPage() int {
	if p.Page < 1 {
		p.Page = pkgConst.DefaultPage
	}
	return p.Page
}

func (p *PaginationRequest) GetLimit() int {
	if p.Limit < 1 {
		p.Limit = pkgConst.DefaultLimit
	}
	return p.Limit
}

func (p *PaginationRequest) GetOffset() int {
	return (p.GetPage() - 1) * p.GetLimit()
}

// GetSearch returns the search keyword.
func (p *PaginationRequest) GetSearch() string {
	return p.Search
}

func (p *PaginationRequest) GetTotal() int64 {
	return p.Total
}

// GetTotalPages calculates total pages from item count and page size.
func (p *PaginationRequest) GetTotalPages() int {
	return int(math.Ceil(float64(p.Total) / float64(p.Limit)))
}

// GetOrderClause returns a SQL ORDER BY clause, e.g. "created_at desc".
func (p *PaginationRequest) GetOrderClause() string {
	sortBy := p.SortBy
	if sortBy == "" {
		sortBy = pkgConst.DefaultOrder
	}
	sortOrder := p.SortOrder
	if sortOrder == "" {
		sortOrder = pkgConst.DefaultSortOrder
	}
	return sortBy + " " + sortOrder
}
