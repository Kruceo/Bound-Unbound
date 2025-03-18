package v1

var ErrorMessages = map[string]string{
	"JSON_ENCODING":        "Failed to encode JSON",
	"JSON_DECODE":          "Failed to decode JSON",
	"AUTH":                 "Unauthorized access",
	"WRONG_LOGIN":          "Invalid login credentials",
	"LOGIN":                "Intern auth error",
	"OVERWRITING_REGISTER": "Intern auth error",
	"BODY_FORMAT":          "Invalid request body format",
	"UNKNOWN_NODE":         "Unknown node identifier",
	"NODE_CONNECTION":      "Failed to connect to node",
}

func GetErrorMessage(code string) string {
	if msg, exists := ErrorMessages[code]; exists {
		return msg
	}
	return "Unknown error"
}
