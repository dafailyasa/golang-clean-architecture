package helpers

import (
	"fmt"
	"math/rand"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var usernameRegex = regexp.MustCompile(`[^a-z0-9]+`)

func IsNumeric(s string) bool {
	if s == "" {
		return false
	}
	_, err := strconv.Atoi(s)
	return err == nil
}

func StructToFormData(v any) (map[string]string, error) {
	rv := reflect.ValueOf(v)

	// Dereference pointer if needed
	if rv.Kind() == reflect.Ptr {
		if rv.IsNil() {
			return nil, fmt.Errorf("StructToFormData: input is a nil pointer")
		}
		rv = rv.Elem()
	}

	if rv.Kind() != reflect.Struct {
		return nil, fmt.Errorf("StructToFormData: expected struct, got %s", rv.Kind())
	}

	result := make(map[string]string)
	rt := rv.Type()

	for i := range rt.NumField() {
		field := rt.Field(i)
		value := rv.Field(i)

		// Skip unexported fields
		if !field.IsExported() {
			continue
		}

		tag := field.Tag.Get("json")
		if tag == "" || tag == "-" {
			continue
		}

		result[tag] = fmt.Sprintf("%v", value.Interface())
	}

	return result, nil
}

func GenerateUsername(firstName, lastName string) string {
	rand.Seed(time.Now().UnixNano())

	first := sanitize(firstName)
	last := sanitize(lastName)

	base := first
	if last != "" {
		base = first + "." + last
	}

	// Add 3-digit random number
	suffix := rand.Intn(900) + 100
	return fmt.Sprintf("%s%d", base, suffix)
}

func sanitize(input string) string {
	input = strings.ToLower(strings.TrimSpace(input))
	return usernameRegex.ReplaceAllString(input, "")
}

func ExtractUUIDLocation(location string) string {
	parts := strings.Split(strings.TrimRight(location, "/"), "/")
	return parts[len(parts)-1]
}
