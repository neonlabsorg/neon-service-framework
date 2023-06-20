package api

type ErrorResponseModel struct {
	//
	// Error
	//
	Error HttpErrorResponseModel `json:"error"`
}

type HttpErrorResponseModel struct {
	//
	// Message
	//
	Message string `json:"message"`
	//
	// Name
	//
	Name string `json:"name"`
	//
	// Status Code
	//
	StatusCode int `json:"status_code"`
	//
	// Context
	//
	Context map[string]string `json:"context"`
}

type SuccessResponse struct {
	Result string `json:"result"`
}

func NewSuccessResponse() *SuccessResponse {
	return &SuccessResponse{
		Result: "success",
	}
}

type ValidationErrorResponseModel struct {
	//
	// Error
	//
	Error ValidationErrorModel `json:"error"`
}

// swagger:model
type ValidationErrorModel struct {
	//
	// Message
	//
	Message string `json:"message"`
	//
	// Name
	//
	Name string `json:"name"`
	//
	// Status Code
	//
	StatusCode int `json:"status_code"`
	//
	// Fields
	//
	Fields []ValidationErrorFieldModel `json:"fields"`
}

// swagger:model
type ValidationErrorFieldModel struct {
	FieldName string `json:"field_name"`
	Namespace string `json:"namespace"`
	Tag       string `json:"tag"`
	TagParam  string `json:"tag_param"`
	Message   string `json:"message"`
}
