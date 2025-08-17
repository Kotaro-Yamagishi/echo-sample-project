package output

import "time"

// ErrorResponse エラーレスポンス用の構造体
type ErrorResponse struct {
	Success   bool      `json:"success"`
	Error     string    `json:"error"`
	Code      int       `json:"code"`
	Timestamp time.Time `json:"timestamp"`
}

// SuccessResponse 成功レスポンス用の構造体
type SuccessResponse struct {
	Success   bool        `json:"success"`
	Data      interface{} `json:"data"`
	Timestamp time.Time   `json:"timestamp"`
}

// NewErrorResponse エラーレスポンスを作成
func NewErrorResponse(code int, message string) ErrorResponse {
	return ErrorResponse{
		Success:   false,
		Error:     message,
		Code:      code,
		Timestamp: time.Now(),
	}
}

// NewSuccessResponse 成功レスポンスを作成
func NewSuccessResponse(data interface{}) SuccessResponse {
	return SuccessResponse{
		Success:   true,
		Data:      data,
		Timestamp: time.Now(),
	}
}
