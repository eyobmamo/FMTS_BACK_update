package utils

import "fmt"

const (
	DefaultPage    = 1
	DefaultPerPage = 10
)

type ContextKey string

// type ErrorDefinition struct {
// 	Code    int    `json:"code"`
// 	Message string `json:"message"`
// }

type Filter struct {
	// page specifies the page number
	Page int `json:"page"`
	// per page specifies the number of results per page
	PerPage int `json:"per_page"`
	// search specifies the search term
	Search string `json:"search"`
	// filters user using status
	Filters string `json:"filters"`
}

func (e ErrorDefinition) Error() string {
	return fmt.Sprintf(`{"code": %d ,"message":"%s"}`, e.Code, e.Message)
}

type SuccesResponse struct {
	Ok         bool `json:"ok"`
	StatusCode int  `json:"status_code"`
	Data       any  `json:"data"`
}
