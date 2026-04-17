package middlewares

import (
	"net/http"
	"os"
	"strings"
)

var (
	allowedOrigins = buildAllowedOrigins()
	allowWildcard  = allowedOrigins["*"]
)

func CORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		if originAllowed(origin) {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Access-Control-Allow-Credentials", "true")
		} else if allowWildcard {
			w.Header().Set("Access-Control-Allow-Origin", "*")
		}

		w.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, X-API-Key, X-localization, Authorization")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, PUT, PATCH, DELETE, OPTIONS")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func buildAllowedOrigins() map[string]bool {
	raw := os.Getenv("APP_ALLOWED_ORIGINS")
	if strings.TrimSpace(raw) == "" {
		raw = "http://localhost:3000"
	}

	origins := make(map[string]bool)
	for _, origin := range strings.Split(raw, ",") {
		trimmed := strings.TrimSpace(origin)
		if trimmed == "" {
			continue
		}
		origins[trimmed] = true
	}

	return origins
}

func originAllowed(origin string) bool {
	if origin == "" {
		return false
	}

	if allowWildcard {
		return true
	}

	return allowedOrigins[origin]
}
