package response

type Response struct {
	Code    int         `json:"code"`
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

type Error struct {
	Code   int    `json:"code"`
	Status string `json:"status"`
	// Message   string      `json:"message"`
	Errors interface{} `json:"errors,omitempty"`
}

// func FailedResponse(message string) map[string]interface{} {
// 	return map[string]interface{}{
// 		"success": false,
// 		"message": message,
// 	}
// }