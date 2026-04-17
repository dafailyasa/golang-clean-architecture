package errors

import (
	"fmt"
	"strings"
)

type ErrorCustomize struct {
	isBusinessError       bool
	isInfrastructureError bool
	codePrefix            string
	codeValue             string
	message               string
}

type ErrorCustomizeDTO struct {
	IsBusinessError       bool   `json:"is_business_error"`
	IsInfrastructureError bool   `json:"is_infrastructure_error"`
	CodePrefix            string `json:"code_prefix"`
	CodeValue             string `json:"code_value"`
	Message               string `json:"message"`
	Code                  string `json:"code"`
}

func (e ErrorCustomize) Error() string {
	return e.message
}

func (e ErrorCustomize) Code() string {
	return fmt.Sprintf("%s%s", strings.ToUpper(e.codePrefix), strings.ToUpper(e.codeValue))
}

func (e ErrorCustomize) CodePrefix() string {
	return e.codePrefix
}

func (e ErrorCustomize) IsBusinessError() bool {
	return e.isBusinessError
}

func (e ErrorCustomize) IsInfrastructureError() bool {
	return e.isInfrastructureError
}

func GetError(err error) ErrorCustomizeDTO {
	if er, ok := err.(*ErrorCustomize); ok {
		return ErrorCustomizeDTO{
			IsBusinessError:       er.isBusinessError,
			IsInfrastructureError: er.isInfrastructureError,
			CodePrefix:            er.codePrefix,
			CodeValue:             er.codeValue,
			Message:               er.message,
			Code:                  er.Code(),
		}
	}

	return ErrorCustomizeDTO{
		IsBusinessError:       false,
		IsInfrastructureError: false,
		CodePrefix:            "",
		CodeValue:             "",
		Message:               err.Error(),
		Code:                  "UNKNOWN_ERROR",
	}
}

func NewBusinessError(codePrefix, codeValue string, message string) error {
	return &ErrorCustomize{
		isBusinessError:       true,
		isInfrastructureError: false,
		codePrefix:            codePrefix,
		codeValue:             codeValue,
		message:               message,
	}
}

func NewTechnicalError(codePrefix, codeValue string, message string) error {
	return &ErrorCustomize{
		isBusinessError:       false,
		isInfrastructureError: false,
		codePrefix:            codePrefix,
		codeValue:             codeValue,
		message:               message,
	}
}

func NewInfrastructureError(codePrefix, codeValue string, message string) error {
	return &ErrorCustomize{
		isBusinessError:       false,
		isInfrastructureError: true,
		codePrefix:            codePrefix,
		codeValue:             codeValue,
		message:               message,
	}
}
