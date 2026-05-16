package middlewares

import (
	pkgErrors "auth-service/pkg/errors"
	pkgResponse "auth-service/pkg/response"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
)

type XUserInfo struct {
	GivenName  string `json:"given_name"`
	FamilyName string `json:"family_name"`
	Email      string `json:"email"`
	Sub        string `json:"sub"`
}

type contextKey string

const UserKey contextKey = "user"

func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		xUserInfo := r.Header.Get("X-Userinfo")
		if xUserInfo == "" {
			pkgResponse.ErrorWithCode(
				w,
				http.StatusUnauthorized,
				"You are not authorized to do this request",
				pkgErrors.NewBusinessError("AUTH", "001", "missing X-Userinfo headers"),
			)
			return
		}

		decodedBytes, err := base64.StdEncoding.DecodeString(xUserInfo)
		if err != nil {
			pkgResponse.ErrorWithCode(
				w,
				http.StatusUnauthorized,
				"You are not authorized to do this request",
				pkgErrors.NewBusinessError("AUTH", "002", fmt.Sprintf("invalid decode X-Userinfo header: %v", err)),
			)
			return
		}

		var user XUserInfo
		if err := json.Unmarshal(decodedBytes, &user); err != nil {
			pkgResponse.ErrorWithCode(
				w,
				http.StatusBadRequest,
				"something went wrong please try again",
				pkgErrors.NewBusinessError("AUTH", "003", fmt.Sprintf("failed to unmarshal: %v", err)),
			)
			return
		}

		ctx := context.WithValue(r.Context(), UserKey, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func GetUserInfoFromContext(ctx context.Context) (XUserInfo, bool) {
	user, ok := ctx.Value(UserKey).(XUserInfo)
	return user, ok
}
