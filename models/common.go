package models

type GenericSuccessResponse struct {
	Result SuccessResult `json:"result"`
}

type GenericErrorResponse struct {
	Result ErrorResult `json:"result"`
}

type SuccessResult struct {
	Status  string `json:"status" example:"Success"`
	Message string `json:"message" example:"Success message"`
}

type ErrorResult struct {
	Status  string `json:"status" example:"Error"`
	Message string `json:"message" example:"Error message"`
}
