package utils

type AppResponse struct {
	Success      bool        `json:"success"`
	Data         interface{} `json:"data"`
	Error        interface{} `json:"error,omitempty"`
	ErrorMessage string      `json:"error_message,omitempty"`
}

func (appResponse *AppResponse) ToMap() map[string]interface{} {
	res := make(map[string]interface{})

	res["success"] = appResponse.Success
	res["data"] = appResponse.Data

	if !appResponse.Success {
		res["error"] = appResponse.Error
		res["error_message"] = appResponse.ErrorMessage
	}

	return res
}

func ARSuccess(data interface{}) AppResponse {
	return AppResponse{
		Success: true,
		Data:    data,
	}
}

func ARFailure(errorMessage string, error ...interface{}) AppResponse {
	if errorMessage == "" {
		errorMessage = "Something got out of the hands."
	}
	return AppResponse{
		Success:      false,
		Data:         nil,
		ErrorMessage: errorMessage,
		Error:        error,
	}
}

func ARValFail(errorMessage string, error []string) AppResponse {
	if errorMessage == "" {
		errorMessage = "Incorrect data was sent."
	}
	return AppResponse{
		Success:      false,
		Data:         nil,
		ErrorMessage: errorMessage,
		Error:        error,
	}
}