package apiresponse

var (
	StatusError   = "error"
	StatusSuccess = "success"
)

type ApiResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

func Error(msg string) ApiResponse {
	return ApiResponse{
		Status:  StatusError,
		Message: msg,
	}
}

func Success(msg string) ApiResponse {
	return ApiResponse{
		Status:  StatusSuccess,
		Message: msg,
	}
}
