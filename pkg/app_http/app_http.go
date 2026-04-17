package app_http

import (
	"auth-service/pkg/constant"
	pkgErrors "auth-service/pkg/errors"
	"auth-service/pkg/strings"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/url"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/valyala/fasthttp"
)

type HttpResponse struct {
	StatusCode int
	Headers    map[string]string
}

type (
	File struct {
		FileName string
		File     io.Reader
	}

	Request struct {
		Method       string
		Endpoint     string
		Headers      map[string]string
		Body         any
		ToURLEncoded bool
		Files        map[string]File
		FormData     map[string]string

		// Optional per-request timeout override (0 = use client default)
		Timeout time.Duration
	}
)

// ClientConfig holds configuration for the HTTP client.
type ClientConfig struct {
	MaxConnsPerHost     int
	MaxIdleConnDuration time.Duration
	ReadTimeout         time.Duration
	WriteTimeout        time.Duration
	MaxRetries          int
	RetryDelay          time.Duration
}

// DefaultClientConfig returns sensible defaults.
func DefaultClientConfig() ClientConfig {
	return ClientConfig{
		MaxConnsPerHost:     50,
		MaxIdleConnDuration: 30 * time.Second,
		ReadTimeout:         30 * time.Second,
		WriteTimeout:        30 * time.Second,
		MaxRetries:          0,
		RetryDelay:          500 * time.Millisecond,
	}
}

type AppHttp struct {
	client *fasthttp.Client
	log    *logrus.Logger
	cfg    ClientConfig
}

// NewClient creates a new AppHttp client with default settings.
func NewClient(log *logrus.Logger) *AppHttp {
	return NewClientWithConfig(log, DefaultClientConfig())
}

// NewClientWithConfig creates a new AppHttp client with custom config.
func NewClientWithConfig(log *logrus.Logger, cfg ClientConfig) *AppHttp {
	return &AppHttp{
		client: &fasthttp.Client{
			MaxConnsPerHost:     cfg.MaxConnsPerHost,
			Dial:                fasthttp.Dial,
			MaxIdleConnDuration: cfg.MaxIdleConnDuration,
			ReadTimeout:         cfg.ReadTimeout,
			WriteTimeout:        cfg.WriteTimeout,
		},
		log: log,
		cfg: cfg,
	}
}

// DoHttpRequest executes an HTTP request and decodes the response into res.
// Supports optional retries, context cancellation, and multipart uploads.
func (h *AppHttp) DoHttpRequest(ctx context.Context, req Request, res any) (*HttpResponse, error) {
	var (
		err error
	)

	var httpResponse *HttpResponse
	maxAttempts := h.cfg.MaxRetries + 1
	for attempt := 1; attempt <= maxAttempts; attempt++ {
		// Check context before each attempt
		if ctx.Err() != nil {
			return nil, pkgErrors.NewTechnicalError("DOREQ", "006", fmt.Sprintf("context cancelled before attempt %d: %s", attempt, ctx.Err()))
		}

		httpResponse, err = h.doOnce(ctx, req, res)
		if err == nil {
			return httpResponse, nil
		}

		if attempt < maxAttempts {
			h.log.WithContext(ctx).Warnf("HTTP request attempt %d/%d failed: %v — retrying in %s", attempt, maxAttempts, err, h.cfg.RetryDelay)
			select {
			case <-ctx.Done():
				return nil, pkgErrors.NewTechnicalError("DOREQ", "006", fmt.Sprintf("context cancelled during retry wait: %s", ctx.Err()))
			case <-time.After(h.cfg.RetryDelay):
			}
		}
	}

	return httpResponse, err
}

// doOnce performs a single HTTP attempt.
// doOnce performs a single HTTP attempt.
func (h *AppHttp) doOnce(ctx context.Context, req Request, res any) (*HttpResponse, error) {
	request := fasthttp.AcquireRequest()
	response := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseRequest(request)
	defer fasthttp.ReleaseResponse(response)

	request.Header.DisableNormalizing()
	request.Header.SetMethod(req.Method)
	request.SetRequestURI(req.Endpoint)

	for key, value := range req.Headers {
		request.Header.Set(key, value)
	}

	if err := h.setRequestBody(request, req); err != nil {
		return nil, err
	}

	// Marshal req.Body → sanitize → return parsed any (map/slice).
	// Passing any to logrus.Fields means it is serialized as a native nested
	// JSON object — no pre-serialized string, no escaped \n characters.
	maskedRequestBody := h.sanitizeRequestBody(req.Body)

	start := time.Now()

	timeout := h.cfg.ReadTimeout
	if req.Timeout > 0 {
		timeout = req.Timeout
	}

	if err := h.client.DoTimeout(request, response, timeout); err != nil {
		h.log.WithContext(ctx).WithFields(logrus.Fields{
			"method":      req.Method,
			"endpoint":    req.Endpoint,
			"duration":    time.Since(start),
			"requestBody": maskedRequestBody,
		}).Errorf("failed to execute HTTP request: %v", err)
		return nil, pkgErrors.NewTechnicalError("DOREQ", "005", fmt.Sprintf("failed to execute HTTP request: %s", err.Error()))
	}

	statusCode := response.StatusCode()
	duration := time.Since(start)
	bodyBytes := response.Body()

	// SanitizeBodyParsed returns any — same reason as above, no \n escaping.
	maskedResponseBody := strings.SanitizeBodyParsed(bodyBytes)
	maskedHeaders := strings.SanitizeHeaders(req.Headers)

	if len(bodyBytes) > 0 {
		if err := json.Unmarshal(bodyBytes, res); err != nil {
			h.log.WithContext(ctx).WithFields(logrus.Fields{
				"method":       req.Method,
				"endpoint":     req.Endpoint,
				"duration":     duration,
				"headers":      maskedHeaders,
				"statusCode":   statusCode,
				"requestBody":  maskedRequestBody,
				"responseBody": maskedResponseBody,
			}).Errorf("failed to decode response: %v", err)

			return nil, pkgErrors.NewTechnicalError(
				"DOREQ",
				"007",
				fmt.Sprintf("failed to decode response: %s", err.Error()),
			)
		}
	}

	headers := make(map[string]string)
	response.Header.VisitAll(func(k, v []byte) {
		headers[string(k)] = string(v)
	})

	h.log.WithContext(ctx).WithFields(logrus.Fields{
		"method":       req.Method,
		"endpoint":     req.Endpoint,
		"duration":     duration,
		"headers":      maskedHeaders,
		"statusCode":   statusCode,
		"requestBody":  maskedRequestBody,
		"responseBody": maskedResponseBody,
	}).Info("success do http request")

	return &HttpResponse{
		StatusCode: statusCode,
		Headers:    headers,
	}, nil
}

// sanitizeRequestBody marshals req.Body → sanitizes sensitive fields → returns
func (h *AppHttp) sanitizeRequestBody(body any) any {
	if body == nil {
		return nil
	}
	raw, err := json.Marshal(body)
	if err != nil {
		return "[unserializable request body]"
	}
	return strings.SanitizeBodyParsed(raw)
}

// setRequestBody sets the body and Content-Type header on the request.
func (h *AppHttp) setRequestBody(request *fasthttp.Request, req Request) error {
	switch {
	case req.FormData != nil && req.ToURLEncoded == true:
		return h.setFormURLEncodedBody(request, req.FormData)
	case req.Files != nil || req.FormData != nil:
		return h.setMultipartBody(request, req)
	case req.Body != nil && req.ToURLEncoded == false:
		return h.setJSONBody(request, req.Body)
	}
	return nil
}

func (h *AppHttp) setMultipartBody(request *fasthttp.Request, req Request) error {
	var buffer bytes.Buffer
	writer := multipart.NewWriter(&buffer)

	for key, file := range req.Files {
		part, err := writer.CreateFormFile(key, file.FileName)
		if err != nil {
			return pkgErrors.NewTechnicalError("DOREQ", "001", fmt.Sprintf("error creating form file: %s", err.Error()))
		}
		if _, err := io.Copy(part, file.File); err != nil {
			return pkgErrors.NewTechnicalError("DOREQ", "002", fmt.Sprintf("could not copy file to form: %s", err.Error()))
		}
	}

	for key, param := range req.FormData {
		if err := writer.WriteField(key, param); err != nil {
			return pkgErrors.NewTechnicalError("DOREQ", "003", fmt.Sprintf("could not write form field %q: %s", key, err.Error()))
		}
	}

	if err := writer.Close(); err != nil {
		return pkgErrors.NewTechnicalError("DOREQ", "004", fmt.Sprintf("could not close multipart writer: %s", err.Error()))
	}

	request.SetBody(buffer.Bytes())
	request.Header.Set("Content-Type", writer.FormDataContentType())
	return nil
}

func (h *AppHttp) setJSONBody(request *fasthttp.Request, body any) error {
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return pkgErrors.NewTechnicalError("DOREQ", "008", fmt.Sprintf("could not marshal request body: %s", err.Error()))
	}
	request.SetBody(jsonBody)
	request.Header.Set("Content-Type", constant.ApplicationJsonConst)
	return nil
}

func (h *AppHttp) setFormURLEncodedBody(request *fasthttp.Request, form map[string]string) error {
	values := url.Values{}
	for key, value := range form {
		values.Set(key, value)
	}
	request.SetBodyString(values.Encode())
	request.Header.Set("Content-Type", constant.FormUrlEncodedConst)
	return nil
}
