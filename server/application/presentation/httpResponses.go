package presentation

import (
	"encoding/json"
	"net/http"
)

type Response[T any] struct {
	Message string `json:"message"`
	Data    T      `json:"data"`
}

type ErrorResponse struct {
	Message   string `json:"message"`
	ErrorCode string `json:"errorCode"`
	Error     bool   `json:"error"`
}

type FastErrorResponses struct {
	errorMessages map[string]string
}

func (f *FastErrorResponses) GetErrorMessage(code string) string {
	if msg, exists := f.errorMessages[code]; exists {
		return msg
	}
	return "Unknown error"
}

func (f *FastErrorResponses) Execute(w http.ResponseWriter, r *http.Request, errorCode string, statusCode int) {
	b := ErrorResponse{Message: f.GetErrorMessage(errorCode), ErrorCode: errorCode, Error: true}
	encoded, err := json.Marshal(b)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(nil)
		return
	}
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(statusCode)
	w.Write(encoded)
}

func NewFastErrorResponses() FastErrorResponses {
	var errorMessages = map[string]string{
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
	return FastErrorResponses{errorMessages: errorMessages}
}
