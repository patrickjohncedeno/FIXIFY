package errors

type ErrorModel struct {
	Message   string      `json:"message"`
	IsSuccess bool        `json:"IsSuccess"`
	Error     interface{} `json:"error"`
}
