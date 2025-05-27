package response

type ResponseModel struct {
	RetCode string      `json:"retCode"` // Fixed the struct tag syntax
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}
