package v1

var ErrorMessages = map[string]string{
	"JSON_ENCODING":        "Failed to encode JSON",
	"JSON_DECODE":          "Failed to decode JSON",
	"AUTH":                 "Unauthorized access",
	"AUTH_BLOCKED":         "Unauthorized access",
	"WRONG_LOGIN":          "Invalid login credentials",
	"LOGIN":                "Intern auth error",
	"OVERWRITING_REGISTER": "Intern auth error",
	"BODY_FORMAT":          "Invalid request body format",
	"UNKNOWN_NODE":         "Unknown node identifier",
	"NODE_CONNECTION":      "Failed to connect to node",
	"NODE_RESPONSE":        "Failed to get response from to node",
	"CONNECTION_SECURITY":  "Failed to get a secure response",
}

func GetErrorMessage(code string) string {
	if msg, exists := ErrorMessages[code]; exists {
		return msg
	}
	return "Unknown error"
}
