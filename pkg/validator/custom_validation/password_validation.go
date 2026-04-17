package custom_validation

import (
	"unicode"

	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
)

var TagPassCustomVal = "password"

func passwordValid(s string) bool {
	if len(s) < 8 || len(s) > 32 {
		return false
	}

	var hasUpper, hasLower, hasNumber, hasSymbol bool

	for _, r := range s {
		switch {
		case unicode.IsUpper(r):
			hasUpper = true
		case unicode.IsLower(r):
			hasLower = true
		case unicode.IsDigit(r):
			hasNumber = true
		case unicode.IsPunct(r) || unicode.IsSymbol(r):
			hasSymbol = true
		}
	}

	return hasUpper && hasLower && hasNumber && hasSymbol
}

func RegisterPasswordValidation(validate *validator.Validate, trans ut.Translator) error {
	// register custom validation
	err := validate.RegisterValidation(TagPassCustomVal, func(fl validator.FieldLevel) bool {
		return passwordValid(fl.Field().String())
	})
	if err != nil {
		return err
	}

	// register custom translation
	err = validate.RegisterTranslation(
		TagPassCustomVal,
		trans,
		func(ut ut.Translator) error {
			return ut.Add(
				TagPassCustomVal,
				"{0} must be 8–32 characters long and include uppercase, lowercase, number, and symbol",
				true,
			)
		},
		func(ut ut.Translator, fe validator.FieldError) string {
			msg, _ := ut.T(TagPassCustomVal, fe.Field())
			return msg
		},
	)
	if err != nil {
		return err
	}

	return nil
}
