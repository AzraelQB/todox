package handlers

import "github.com/go-playground/validator/v10"

type Response struct {
	Todo TodoResponse `json:"todo"`
}

type TodosResponse struct {
	Todos []TodoResponse `json:"todos"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type SuccessResponse struct {
	Message string `json:"message"`
}

var validate = validator.New(validator.WithRequiredStructEnabled())
