// Package context contains utility methods for context related functionalities.
package contexts

import (
	constant "FMTS/utils"
	"context"
	"net/http"
)

type UserContext struct {
	UserCode    string
	UserID      string
	FullName    string
	PhoneNumber string
	UserRole    string
}

func ExtractUserContext(r *http.Request) UserContext {
	// This method extracts users data from the middleware context
	get := func(key string) string {
		val, _ := r.Context().Value(constant.ContextKey(key)).(string)
		return val
	}

	return UserContext{
		// UserCode:    get("user_code"),
		UserID:      get("user_id"),
		FullName:    get("full_name"),
		PhoneNumber: get("phone_number"),
		UserRole:    get("user_role"),
	}
}

func ExtractContext(c context.Context) UserContext {
	// This method extracts users data from the middleware context
	get := func(key string) string {
		val, _ := c.Value(constant.ContextKey(key)).(string)
		return val
	}
	return UserContext{
		UserCode:    get("user_code"),
		UserID:      get("user_id"),
		FullName:    get("full_name"),
		PhoneNumber: get("phone_number"),
		UserRole:    get("user_role"),
	}
}

func (u UserContext) IsIncomplete() bool {
	return u.UserID == "" || u.FullName == "" || u.PhoneNumber == ""
}
