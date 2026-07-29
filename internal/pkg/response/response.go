package response

import (
	"encoding/json"
	"net/http"
)

type APIError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}
type Meta struct {
	Total     int `json:"total"`
	Page      int `json:"page"`
	PerPage   int `json:"per_page"`
	TotalPage int `json:"total_page"`
}

type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   *APIError   `json:"error,omitempty"`
	Meta    *Meta       `json:"meta,omitempty"`
}

func jsonResponse(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if data != nil {
		_ = json.NewEncoder(w).Encode(data)
	}
}

// 200 OK
func OK(w http.ResponseWriter, message string, data interface{}) {
	jsonResponse(w, http.StatusOK, Response{
		Success: true,
		Message: message,
		Data:    data,
	})
}

// 200 OK with Meta (Pagination, etc...)
func OKWithMeta(w http.ResponseWriter, message string, data interface{}, meta *Meta) {
	jsonResponse(w, http.StatusOK, Response{
		Success: true,
		Message: message,
		Data:    data,
		Meta:    meta,
	})
}

// 201 Created
func Created(w http.ResponseWriter, message string, data interface{}) {
	jsonResponse(w, http.StatusCreated, Response{
		Success: true,
		Message: message,
		Data:    data,
	})
}

// 204 No Content
func NoContent(w http.ResponseWriter) {
	jsonResponse(w, http.StatusNoContent, Response{
		Success: true,
	})
}

// 400 Bad Request
func BadRequest(w http.ResponseWriter, code, message, details string) {
	jsonResponse(w, http.StatusBadRequest, Response{
		Success: false,
		Message: message,
		Error: &APIError{
			Code:    code,
			Message: message,
			Details: details,
		},
	})
}

// 401 Unauthorized
func Unauthorized(w http.ResponseWriter, message string) {
	jsonResponse(w, http.StatusUnauthorized, Response{
		Success: false,
		Message: message,
		Error: &APIError{
			Code:    "UNAUTHORIZED",
			Message: message,
		},
	})
}

// 403 Forbidden
func Forbidden(w http.ResponseWriter, message string) {
	jsonResponse(w, http.StatusForbidden, Response{
		Success: false,
		Message: message,
		Error: &APIError{
			Code:    "FORBIDDEN",
			Message: message,
		},
	})
}

// 404 Not Found
func NotFound(w http.ResponseWriter, message string) {
	jsonResponse(w, http.StatusNotFound, Response{
		Success: false,
		Message: message,
		Error: &APIError{
			Code:    "NOT_FOUND",
			Message: message,
		},
	})
}

// 405 Method Not Allowed
func MethodNotAllowed(w http.ResponseWriter, message string) {
	jsonResponse(w, http.StatusMethodNotAllowed, Response{
		Success: false,
		Message: message,
		Error: &APIError{
			Code:    "METHOD_NOT_ALLOWED",
			Message: message,
		},
	})
}

// 409 Conflict
func Conflict(w http.ResponseWriter, message string) {
	jsonResponse(w, http.StatusConflict, Response{
		Success: false,
		Message: message,
		Error: &APIError{
			Code:    "CONFLICT",
			Message: message,
		},
	})
}

// 429 Too Many Requests
func TooManyRequests(w http.ResponseWriter, message string) {
	jsonResponse(w, http.StatusTooManyRequests, Response{
		Success: false,
		Message: message,
		Error: &APIError{
			Code:    "TOO_MANY_REQUESTS",
			Message: message,
		},
	})
}

// 500 Internal Server Error
func InternalServerError(w http.ResponseWriter, message string) {
	jsonResponse(w, http.StatusInternalServerError, Response{
		Success: false,
		Message: message,
		Error: &APIError{
			Code:    "INTERNAL_SERVER_ERROR",
			Message: message,
		},
	})
}

// 503 Service Unavailable
func ServiceUnavailable(w http.ResponseWriter, message string) {
	jsonResponse(w, http.StatusServiceUnavailable, Response{
		Success: false,
		Message: message,
		Error: &APIError{
			Code:    "SERVICE_UNAVAILABLE",
			Message: message,
		},
	})
}
