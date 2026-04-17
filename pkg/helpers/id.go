package helpers

import (
	"strconv"

	pkgErrors "auth-service/pkg/errors"
)

// ParseUintParam converts a decimal string ID into uint and wraps validation errors
// with a consistent business error code.
func ParseUintParam(param, codePrefix, codeValue, message string) (uint, error) {
	parsed, err := strconv.ParseUint(param, 10, 0)
	if err != nil {
		return 0, pkgErrors.NewBusinessError(codePrefix, codeValue, message)
	}
	return uint(parsed), nil
}
