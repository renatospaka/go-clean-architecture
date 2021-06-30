package utils

//service error should be returned when business error occur
type ServiceError struct {
	Message string `json:"message"`
} 