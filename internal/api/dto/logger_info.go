package dto

type LoggerInfo struct {
	Scope   string `json:"category"`
	Message string `json:"message"`
}
