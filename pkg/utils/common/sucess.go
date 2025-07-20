package common

type SuccessDefinition struct {
	Code    string
	Message string
}

type SuccessGroup map[string]SuccessDefinition

type SuccessDefinitions struct {
	General SuccessGroup
	Auth    SuccessGroup
	User    SuccessGroup
}

var DefineSuccess = SuccessDefinitions{
	General: SuccessGroup{
		"CONFLICT_KEY": {
			Code:    "GEN_001",
			Message: "Duplicate key success: The specified field already exists.",
		},
		"SUCCESS": {
			Code:    "GEN_002",
			Message: "Success: The operation was successful.",
		},
	},
	Auth: SuccessGroup{
		"AUTH_USER_OTP_SENT_Customer": {
			Code:    "AUTH_SUCCESS_001",
			Message: "Hello Customer your are logged in for the first time it needs to verified the OTP is sent to your phone",
		},
	},
}

func GetSuccessResponseByCode(code string) (SuccessDefinition, bool) {
	for _, group := range []SuccessGroup{
		DefineSuccess.General,
		DefineSuccess.Auth,
		DefineSuccess.User,
	} {
		for _, def := range group {
			if def.Code == code {
				return def, true
			}
		}
	}
	return SuccessDefinition{}, false
}
