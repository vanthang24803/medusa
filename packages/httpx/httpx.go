package httpx

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"reflect"
	"strconv"
	"time"

	"ecommerce/packages/types"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/google/uuid"
)

type contextKey string

// Metadata holds response metadata.
type Metadata struct {
	RequestID string `json:"requestId"`
	Timestamp string `json:"timestamp"`
	LatencyUs int64  `json:"latencyUs"`
}

const startTimeKey contextKey = "start_time"

// GenerateUUID generates a compliant UUID v4.
func GenerateUUID() string {
	id, err := uuid.NewV7()
	if err != nil {
		panic(err)
	}
	return id.String()
}

// RequestIDUUIDMiddleware ensures X-Request-Id is populated with a UUID v4 if not already present.
func RequestIDUUIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-Request-Id") == "" {
			r.Header.Set("X-Request-Id", GenerateUUID())
		}
		next.ServeHTTP(w, r)
	})
}

// RecordStartTimeMiddleware records the start time of the request.
func RecordStartTimeMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), startTimeKey, time.Now())
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func getMetadata(r *http.Request) Metadata {
	var latencyUs int64
	if r != nil {
		if startTime, ok := r.Context().Value(startTimeKey).(time.Time); ok {
			latencyUs = time.Since(startTime).Microseconds()
		}
	}

	var reqID string
	if r != nil {
		reqID = middleware.GetReqID(r.Context())
	}

	return Metadata{
		RequestID: reqID,
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		LatencyUs: latencyUs,
	}
}

// JSON writes a JSON response with status code and metadata at the end of the object.
func JSON(w http.ResponseWriter, r *http.Request, status int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	metadata := getMetadata(r)
	metaBytes, err := json.Marshal(metadata)
	if err != nil {
		res := map[string]any{
			"error": map[string]string{"code": "internal_error", "message": err.Error()},
		}
		_ = json.NewEncoder(w).Encode(res)
		return
	}

	if body == nil {
		_, _ = w.Write([]byte(`{"metadata":` + string(metaBytes) + `}`))
		return
	}

	bodyBytes, err := json.Marshal(body)
	if err != nil {
		res := map[string]any{
			"error":    map[string]string{"code": "internal_error", "message": err.Error()},
			"metadata": metadata,
		}
		_ = json.NewEncoder(w).Encode(res)
		return
	}

	// Trim whitespace
	trimmed := bytes.TrimSpace(bodyBytes)
	if len(trimmed) >= 2 && trimmed[0] == '{' && trimmed[len(trimmed)-1] == '}' {
		// Trim the trailing '}'
		stripped := trimmed[:len(trimmed)-1]

		// If the original object was not empty, append a comma separator
		var result []byte
		if len(bytes.TrimSpace(stripped)) > 1 {
			result = append(stripped, ',')
		} else {
			result = stripped
		}

		result = append(result, []byte(`"metadata":`)...)
		result = append(result, metaBytes...)
		result = append(result, '}')

		_, _ = w.Write(result)
		return
	}

	// If it is not a JSON object (like slice, string, number), wrap it under data
	res := map[string]any{
		"data":     body,
		"metadata": metadata,
	}
	_ = json.NewEncoder(w).Encode(res)
}

// Error maps AppError / sentinel error to the appropriate HTTP status.
func Error(w http.ResponseWriter, r *http.Request, err error) {
	if appErr, ok := errors.AsType[*types.AppError](err); ok {
		JSON(w, r, statusForCode(appErr.Code), map[string]any{
			"error": map[string]string{"code": appErr.Code, "message": appErr.Message},
		})
		return
	}
	switch {
	case errors.Is(err, types.ErrNotFound):
		JSON(w, r, http.StatusNotFound, errBody("not_found", err.Error()))
	case errors.Is(err, types.ErrConflict):
		JSON(w, r, http.StatusConflict, errBody("conflict", err.Error()))
	case errors.Is(err, types.ErrValidation):
		JSON(w, r, http.StatusBadRequest, errBody("validation_error", err.Error()))
	case errors.Is(err, types.ErrUnauthorized):
		JSON(w, r, http.StatusUnauthorized, errBody("unauthorized", err.Error()))
	case errors.Is(err, types.ErrForbidden):
		JSON(w, r, http.StatusForbidden, errBody("forbidden", err.Error()))
	default:
		JSON(w, r, http.StatusInternalServerError, errBody("internal_error", "internal server error"))
	}
}

func errBody(code, msg string) map[string]any {
	return map[string]any{"error": map[string]string{"code": code, "message": msg}}
}

func statusForCode(code string) int {
	switch code {
	case "not_found":
		return http.StatusNotFound
	case "conflict":
		return http.StatusConflict
	case "validation_error":
		return http.StatusBadRequest
	case "unauthorized":
		return http.StatusUnauthorized
	case "forbidden":
		return http.StatusForbidden
	default:
		return http.StatusInternalServerError
	}
}

// DecodeJSON parses request body into dst, returning a validation error on failure.
func DecodeJSON(r *http.Request, dst any) error {
	if err := json.NewDecoder(r.Body).Decode(dst); err != nil {
		return types.NewValidation("invalid JSON body")
	}
	return nil
}

// DecodeQuery parses URL query params into dst using reflection.
// Supports form tags on struct fields (e.g. `form:"page"`).
func DecodeQuery(r *http.Request, dst any) error {
	values := r.URL.Query()
	return decodeValues(values, dst)
}

func Response(name string, data any) map[string]any {
	return map[string]any{name: data}
}

func decodeValues(values url.Values, dst any) error {
	rv := reflect.ValueOf(dst)
	if rv.Kind() != reflect.Ptr || rv.Elem().Kind() != reflect.Struct {
		return types.NewValidation("decode target must be a pointer to struct")
	}
	rv = rv.Elem()
	rt := rv.Type()

	for i := 0; i < rt.NumField(); i++ {
		field := rt.Field(i)
		fieldVal := rv.Field(i)

		tag := field.Tag.Get("form")
		if tag == "" || tag == "-" {
			continue
		}

		raw, ok := values[tag]
		if !ok || len(raw) == 0 {
			continue
		}
		val := raw[0]

		switch fieldVal.Kind() {
		case reflect.String:
			fieldVal.SetString(val)
		case reflect.Int, reflect.Int64:
			if n, err := strconv.ParseInt(val, 10, 64); err == nil {
				fieldVal.SetInt(n)
			}
		case reflect.Float64:
			if f, err := strconv.ParseFloat(val, 64); err == nil {
				fieldVal.SetFloat(f)
			}
		case reflect.Bool:
			if b, err := strconv.ParseBool(val); err == nil {
				fieldVal.SetBool(b)
			}
		case reflect.Pointer:
			switch fieldVal.Type().Elem().Kind() {
			case reflect.String:
				fieldVal.Set(reflect.ValueOf(&val))
			case reflect.Int:
				if n, err := strconv.Atoi(val); err == nil {
					fieldVal.Set(reflect.ValueOf(&n))
				}
			case reflect.Int64:
				if n, err := strconv.ParseInt(val, 10, 64); err == nil {
					fieldVal.Set(reflect.ValueOf(&n))
				}
			}
		}
	}
	return nil
}
