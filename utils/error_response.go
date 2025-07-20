package utils

import (
	"encoding/json"
	"fmt"
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type ErrorDefinition struct {
	Code    string
	Message string
}

type APIErrorResponse struct {
	Message       string         `json:"message"`
	Code          string         `json:"code,omitempty"`
	AccountLocked bool           `json:"accountLocked,omitempty"`
	Auth          bool           `json:"auth,omitempty"`
	Errors        map[string]any `json:"errors,omitempty"`
}

// ---- Central error definitions ----

var GeneralErrors = map[string]ErrorDefinition{
	"INVALID_ID":  {Code: "INVALID_ID", Message: "Invalid ID format."},
	"NOT_FOUND":   {Code: "NOT_FOUND", Message: "Resource not found."},
	"GEN_UNKNOWN": {Code: "GEN_UNKNOWN", Message: "Unknown error occurred."},
}

var AuthErrors = map[string]ErrorDefinition{
	"AUTH_USER_NOT_FOUND": {Code: "AUTH_USER_NOT_FOUND", Message: "User not found."},
	"INVALID_PASSWORD":    {Code: "INVALID_PASSWORD", Message: "Password is invalid."},
}

// Add more groups here as needed

var AllErrorGroups = []map[string]ErrorDefinition{
	GeneralErrors,
	AuthErrors,
}

// ---- Status mapping (reduced) ----

var errorKeyToStatus = map[string]int{
	"INVALID_ID":          http.StatusBadRequest,
	"NOT_FOUND":           http.StatusNotFound,
	"AUTH_USER_NOT_FOUND": http.StatusNotFound,
	"INVALID_PASSWORD":    http.StatusUnauthorized,
}

// ---- Error lookup ----

func lookupErrorDefinition(key string) ErrorDefinition {
	for _, group := range AllErrorGroups {
		if def, ok := group[key]; ok {
			return def
		}
	}
	return ErrorDefinition{Code: key, Message: key}
}

func statusForErrorKey(key string, provided int, isValidation bool) int {
	if isValidation {
		return http.StatusBadRequest
	}
	if provided != 0 {
		return provided
	}
	if code, ok := errorKeyToStatus[key]; ok {
		return code
	}
	return http.StatusInternalServerError
}

// ---- Validation formatting ----

func formatValidationErrors(ve validation.Errors) map[string]string {
	result := make(map[string]string)
	for field, err := range ve {
		def := lookupErrorDefinition(err.Error())
		result[field] = def.Message
	}
	return result
}

// ---- Error response writer ----

func SendErrorResponse(w http.ResponseWriter, err any, statusOverride int, extra map[string]any) {
	var (
		code              string
		message           string
		errorKey          string
		errors            map[string]any
		isValidationError bool
	)

	switch v := err.(type) {
	case string:
		def := lookupErrorDefinition(v)
		code = def.Code
		message = def.Message
		errorKey = def.Code

	case error:
		if ve, ok := v.(validation.Errors); ok {
			code = "VALIDATION_ERROR"
			message = "One or more required fields are invalid."
			// Convert map[string]string to map[string]any
			formatted := formatValidationErrors(ve)
			errors = make(map[string]any, len(formatted))
			for k, v := range formatted {
				errors[k] = v
			}
			errorKey = code
			isValidationError = true
		} else {
			def := lookupErrorDefinition(v.Error())
			code = def.Code
			message = def.Message
			errorKey = def.Code
		}

	default:
		code = "GEN_UNKNOWN"
		message = fmt.Sprintf("%v", v)
		errorKey = "GEN_UNKNOWN"
	}

	status := statusForErrorKey(errorKey, statusOverride, isValidationError)

	resp := APIErrorResponse{
		Message: message,
		Code:    code,
	}

	if errors != nil {
		resp.Errors = errors
	}

	// Attach extra fields if provided
	if extra != nil {
		if v, ok := extra["accountLocked"].(bool); ok {
			resp.AccountLocked = v
		}
		if v, ok := extra["auth"].(bool); ok {
			resp.Auth = v
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		fmt.Printf("Failed to encode error response: %v\n", err)
	}
}

// Optional shortcut
func HandleServiceError(w http.ResponseWriter, err error) {
	SendErrorResponse(w, err, 0, nil)
}
