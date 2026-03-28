package presenters

type ErrorResponse struct {
	StatusCode   int    `json:"status_code"`
	RequestID    string `json:"request_id"`
	ErrorMessage string `json:"error_message"`
	ErrorType    string `json:"error_type"`
}

type SuccessResponse struct {
	StatusCode int         `json:"status_code"`
	RequestID  string      `json:"request_id"`
	Message    string      `json:"message"`
	Data       interface{} `json:"data"`
}

// func UserErrorResponse(err error) *fiber.Map {
// 	return &fiber.Map{
// 		"status": "fail",
// 		"data":   "",
// 		"error":  err.Error(),
// 	}
// }
