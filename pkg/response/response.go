package response

import (
	"auth-service/pkg/constant"
	"auth-service/pkg/exception"
	pkgPagination "auth-service/pkg/pagination"
	"encoding/json"
	"net/http"
)

type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
	Error   *ErrorInfo  `json:"error,omitempty"`
	Meta    *Meta       `json:"meta,omitempty"`
}

type ErrorInfo struct {
	Code                  string `json:"code"`
	Message               string `json:"message"`
	IsBusinessError       bool   `json:"isBusinessError,omitempty"`
	IsInfrastructureError bool   `json:"isInfrastructureError,omitempty"`
}

type Meta struct {
	Page       int   `json:"page,omitempty"`
	Limit      int   `json:"limit,omitempty"`
	TotalItems int64 `json:"totalItems"`
	TotalPages int   `json:"totalPages"`
}

func writeJSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", constant.ApplicationJsonConst)
	w.WriteHeader(status)

	_ = json.NewEncoder(w).Encode(payload)
}

// Success creates a successful response
func Success(w http.ResponseWriter, message string, data interface{}) {
	writeJSON(w, http.StatusOK, Response{
		Success: true,
		Message: message,
		Data:    data,
	})
}

// SuccessWithMeta creates a successful response with metadata
func SuccessWithMetaPagination(w http.ResponseWriter, message string, data interface{}, pagination *pkgPagination.PaginationRequest) {
	writeJSON(w, http.StatusOK, Response{
		Success: true,
		Message: message,
		Data:    data,
		Meta: &Meta{
			Page:       pagination.GetPage(),
			Limit:      pagination.GetLimit(),
			TotalItems: pagination.GetTotal(),
			TotalPages: pagination.GetTotalPages(),
		},
	})
}

// Created creates a 201 Created response
func Created(w http.ResponseWriter, message string, data interface{}) {
	writeJSON(w, http.StatusCreated, Response{
		Success: true,
		Message: message,
		Data:    data,
	})
}

func SuccessWithNoContent(w http.ResponseWriter, message string) {
	writeJSON(w, http.StatusNoContent, Response{
		Success: true,
		Message: message,
		Data:    nil,
	})
}

// ErrorWithCode creates an error response with custom status code
func ErrorWithCode(w http.ResponseWriter, statusCode int, msg string, err error) {
	errDTO := exception.GetError(err)

	writeJSON(w, statusCode, Response{
		Success: false,
		Message: msg,
		Error: &ErrorInfo{
			Code:                  errDTO.Code,
			Message:               errDTO.Message,
			IsBusinessError:       errDTO.IsBusinessError,
			IsInfrastructureError: errDTO.IsInfrastructureError,
		},
	})
}

func Error(w http.ResponseWriter, msg string, err error) {
	errDTO := exception.GetError(err)
	statusCode := http.StatusBadRequest
	if status := exception.GetHTTPStatus(err); status != nil {
		statusCode = *status
	} else if errDTO.IsBusinessError {
		statusCode = http.StatusUnprocessableEntity
	} else if errDTO.IsInfrastructureError {
		statusCode = http.StatusInternalServerError
	}

	data := Response{
		Success: false,
		Message: msg,
		Error: &ErrorInfo{
			Code:                  errDTO.Code,
			Message:               errDTO.Message,
			IsBusinessError:       errDTO.IsBusinessError,
			IsInfrastructureError: errDTO.IsInfrastructureError,
		},
	}

	writeJSON(w, statusCode, data)
}
