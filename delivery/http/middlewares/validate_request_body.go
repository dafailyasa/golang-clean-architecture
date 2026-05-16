package middlewares

import (
	pkgErrors "auth-service/pkg/errors"
	pkgResponse "auth-service/pkg/response"
	"context"
	"encoding/json"
	"net/http"

	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
)

type ctxKey string

const BodyKey ctxKey = "validatorBody"

func ValidateRequestBody[T any](
	v *validator.Validate,
	trans ut.Translator,
) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var body T

			decoder := json.NewDecoder(r.Body)
			decoder.DisallowUnknownFields()

			if err := decoder.Decode(&body); err != nil {
				pkgResponse.ErrorWithCode(
					w,
					http.StatusBadRequest,
					"unsupported media type",
					pkgErrors.NewBusinessError("UNSUPPORTED_MEDIA_TYPE", "BODY_REQUEST", "Content-Type must be application/json"),
				)
				return
			}

			if err := v.Struct(&body); err != nil {
				ve := err.(validator.ValidationErrors)[0]
				msg := ve.Translate(trans)

				pkgResponse.ErrorWithCode(
					w,
					http.StatusUnprocessableEntity,
					"invalid body format request",
					pkgErrors.NewBusinessError("INVALID_", "BODY_FORMAT_REQUEST", msg),
				)
				return
			}

			ctx := context.WithValue(r.Context(), BodyKey, body)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
