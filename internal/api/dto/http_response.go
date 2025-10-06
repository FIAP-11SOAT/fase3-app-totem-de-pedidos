package dto

type HttpResponse struct {
	Message string `json:"message"`
}

type HttpResponseError struct {
	Error string `json:"error"`
}
