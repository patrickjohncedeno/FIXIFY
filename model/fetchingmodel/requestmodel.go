package fetchingmodel

import "time"

type Request struct {
	Id             int       `json:"request_id"`
	UserId         int       `json:"user_id"`
	FixerId        int       `json:"fixer_id"`
	CategoryId     string    `json:"category_id"`
	Description    string    `json:"description"`
	Status         string    `json:"status"`
	RequestDate    time.Time `json:"request_date"`
	CompletionDate time.Time `json:"completion_date"`
}

// TableName method to specify the correct table name
func (Request) TableName() string {
	return "service_requests" // Explicitly define the table name here
}
