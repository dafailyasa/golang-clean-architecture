package validator

import (
	pkgCustValidation "auth-service/pkg/validator/custom_validation"

	en "github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	enTranslations "github.com/go-playground/validator/v10/translations/en"
)

func NewValidator() (*validator.Validate, ut.Translator, error) {
	validate := validator.New()

	english := en.New()
	uni := ut.New(english, english)

	trans, _ := uni.GetTranslator("en")

	// custom validation
	if err := pkgCustValidation.RegisterPasswordValidation(validate, trans); err != nil {
		return nil, nil, err
	}
	if err := pkgCustValidation.RegisterJWTValidation(validate, trans); err != nil {
		return nil, nil, err
	}

	// translations
	if err := enTranslations.RegisterDefaultTranslations(validate, trans); err != nil {
		return nil, nil, err
	}

	return validate, trans, nil
}
