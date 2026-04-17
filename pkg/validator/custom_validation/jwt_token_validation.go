package custom_validation

import (
	"strings"

	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
)

var TagJWT = "jwtToken"

func RegisterJWTValidation(validate *validator.Validate, trans ut.Translator) error {
	// register custom validation
	err := validate.RegisterValidation(TagJWT, func(fl validator.FieldLevel) bool {
		token := fl.Field().String()
		parts := strings.Split(token, ".")
		return len(parts) == 3 && parts[0] != "" && parts[1] != "" && parts[2] != ""
	})
	if err != nil {
		return err
	}

	// register custom translation
	err = validate.RegisterTranslation(
		TagJWT,
		trans,
		func(ut ut.Translator) error {
			return ut.Add(
				TagJWT,
				"{0} must be a valid JWT token",
				true,
			)
		},
		func(ut ut.Translator, fe validator.FieldError) string {
			msg, _ := ut.T(TagJWT, fe.Field())
			return msg
		},
	)
	if err != nil {
		return err
	}

	return nil
}
