package users

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
)

// Custom time types
type TimeWithoutTimezone time.Time
type TimeWithDate time.Time
type TimeOnly time.Time

// Implement Scanner and Valuer interfaces for TimeWithoutTimezone
func (t *TimeWithoutTimezone) Scan(value interface{}) error {
	if value == nil {
		*t = TimeWithoutTimezone(time.Time{})
		return nil
	}

	switch v := value.(type) {
	case string:
		parsedTime, err := time.Parse("15:04:05", v)
		if err != nil {
			return err
		}
		*t = TimeWithoutTimezone(parsedTime)
	case time.Time:
		*t = TimeWithoutTimezone(v)
	default:
		return fmt.Errorf("unsupported time format %T", v)
	}
	return nil
}

func (t TimeWithoutTimezone) Value() (driver.Value, error) {
	return time.Time(t).Format("15:04:05"), nil
}

// Implement Scanner and Valuer interfaces for TimeWithDate
func (t *TimeWithDate) Scan(value interface{}) error {
	if value == nil {
		*t = TimeWithDate(time.Time{})
		return nil
	}

	switch v := value.(type) {
	case string:
		formats := []string{
			"2006-01-02 15:04:05",
			time.RFC3339,
			"2006-01-02",
			"15:04:05",
			time.RFC3339Nano,
		}

		var parsedTime time.Time
		var err error

		for _, format := range formats {
			parsedTime, err = time.Parse(format, v)
			if err == nil {
				*t = TimeWithDate(parsedTime)
				return nil
			}
		}
		return fmt.Errorf("failed to parse time '%s': %v", v, err)

	case time.Time:
		*t = TimeWithDate(v)
		return nil

	default:
		return fmt.Errorf("unsupported time format %T", v)
	}
}

func (t TimeWithDate) Value() (driver.Value, error) {
	return time.Time(t), nil
}

func (t TimeWithDate) MarshalJSON() ([]byte, error) {
	return json.Marshal(time.Time(t).Format("2006-01-02 15:04:05"))
}

func (t *TimeWithDate) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	formats := []string{
		"2006-01-02 15:04:05",
		"2006-01-02",
		time.RFC3339,
	}

	var parsedTime time.Time
	var err error

	for _, format := range formats {
		parsedTime, err = time.Parse(format, s)
		if err == nil {
			*t = TimeWithDate(parsedTime)
			return nil
		}
	}
	return fmt.Errorf("failed to parse JSON time '%s': %v", s, err)
}

// Implement Scanner and Valuer interfaces for TimeOnly
func (t *TimeOnly) Scan(value interface{}) error {
	if value == nil {
		*t = TimeOnly(time.Time{})
		return nil
	}

	switch v := value.(type) {
	case string:
		parsedTime, err := time.Parse("15:04:05", v)
		if err != nil {
			return err
		}
		*t = TimeOnly(parsedTime)
	case time.Time:
		*t = TimeOnly(v)
	default:
		return fmt.Errorf("unsupported time format %T", v)
	}
	return nil
}

func (t TimeOnly) Value() (driver.Value, error) {
	return time.Time(t).Format("15:04:05"), nil
}

// User model remains the same
type User struct {
	UserId          uint      `gorm:"primaryKey;column:user_id" json:"user_id"`
	Type            string    `json:"type" gorm:"default:'Client'"`
	Gender          string    `gorm:"column:gender" json:"gender"`
	First_name      string    `gorm:"column:first_name" json:"first_name"`
	Last_name       string    `gorm:"column:last_name" json:"last_name"`
	Email           string    `gorm:"column:email;unique" json:"email"`
	Password        string    `gorm:"column:password" json:"password"`
	Phone           string    `gorm:"column:phone;unique" json:"phone"`
	Profile_picture []byte    `gorm:"column:profile_picture;type:bytea" json:"profile_picture"`
	Availability    string    `gorm:"column:availability" json:"availability"`
	Address         string    `gorm:"column:address" json:"address"`
	Created_at      time.Time `gorm:"column:createdat;autoCreateTime" json:"created_at"`
	Updated_at      time.Time `gorm:"column:updatedat;autoUpdateTime" json:"updated_at"`
	FCMToken        string    `gorm:"size:255" json:"fcm_token"`
}

// Repairman model remains the same
type Repairman struct {
	UserId              uint      `gorm:"primaryKey;column:user_id" json:"user_id"`
	Type                string    `json:"type" gorm:"default:'Repairman'"`
	Gender              string    `gorm:"column:gender" json:"gender"`
	First_name          string    `gorm:"column:first_name" json:"first_name"`
	Last_name           string    `gorm:"column:last_name" json:"last_name"`
	Email               string    `gorm:"column:email;unique" json:"email"`
	Password            string    `gorm:"column:password" json:"password"`
	Phone               string    `gorm:"column:phone;unique" json:"phone"`
	Address             string    `gorm:"colßumn:address" json:"address"`
	Profile_picture     []byte    `gorm:"column:profile_picture;type:bytea" json:"profile_picture"`
	Availability        string    `gorm:"column:availability" json:"availability"`
	Verification_status string    `gorm:"column:verification_status" json:"verification_status"`
	Average_rating      float32   `gorm:"column:average_rating" json:"average_rating"`
	Created_at          time.Time `gorm:"column:createdat;autoCreateTime" json:"created_at"`
	Updated_at          time.Time `gorm:"column:updatedat;autoUpdateTime" json:"updated_at"`
	CategoryId          int       `gorm:"column:category_id" json:"category_id"`

	ServiceCategory ServiceCategory `gorm:"foreignKey:CategoryId;references:category_id" json:"service_category"`
}

// Updated EmailVer with TimeWithoutTimezone
type EmailVer struct {
	Email     string              `gorm:"column:email" json:"email"`
	Code      string              `gorm:"column:code" json:"code"`
	Createdat TimeWithoutTimezone `gorm:"column:createdat;autoCreateTime" json:"createdat"`
}

// ServiceCategory remains the same

type ServiceCategory struct {
	CategoryId   int       `gorm:"primaryKey;column:category_id;autoIncrement" json:"category_id"`
	CategoryName string    `gorm:"column:category_name" json:"category_name"`
	Description  string    `gorm:"column:description" json:"description"`
	Created_at   time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	Is_active    bool      `gorm:"column:is_active" json:"is_active"`
}
type ServiceRequest struct {
	RequestId      int          `gorm:"primaryKey;column:request_id;autoIncrement" json:"request_id"`
	UserId         uint         `gorm:"column:user_id" json:"user_id"`
	RepairmanId    uint         `gorm:"column:fixer_id" json:"fixer_id"`
	CategoryId     int          `gorm:"column:category_id" json:"category_id"`
	Status         string       `gorm:"column:status" json:"status"`
	Description    string       `gorm:"column:description" json:"description"`
	RequestDate    TimeWithDate `gorm:"column:request_date;type:timestamp;autoCreateTime" json:"request_date"`
	CompletionDate TimeWithDate `gorm:"column:completion_date;type:timestamp;autoUpdateTime" json:"completion_date"`
	ReviewId       uint         `gorm:"column:review_id" json:"review_id"`

	User            User            `gorm:"foreignKey:UserId;references:UserId" json:"user"`
	Repairman       Repairman       `gorm:"foreignKey:RepairmanId;references:UserId" json:"repairman"` // fixed
	ServiceCategory ServiceCategory `gorm:"foreignKey:CategoryId;references:CategoryId" json:"service_category"`
	Review          *Review         `gorm:"foreignKey:ReviewId;references:ReviewId" json:"review"`
}

type Review struct {
	ReviewId    int          `gorm:"primaryKey;column:review_id" json:"review_id"`
	RequestId   uint         `gorm:"column:request_id" json:"request_id"`
	ClientId    uint         `gorm:"column:client_id" json:"client_id"`
	RepairmanId uint         `gorm:"column:repairman_id" json:"repairman_id"`
	Rating      float64      `gorm:"column:rating" json:"rating"`
	ReviewText  string       `gorm:"column:review_text" json:"review_text"`
	ReviewDate  TimeWithDate `gorm:"column:review_date" json:"review_date"`

	Client    User            `gorm:"foreignKey:ClientId;references:UserId" json:"client"`
	Repairman Repairman       `gorm:"foreignKey:RepairmanId;references:UserId" json:"repairman"`
	Request   *ServiceRequest `gorm:"foreignKey:RequestId;references:RequestId" json:"request"` // 🛠 Make it a pointer
}

type UserVerification struct {
	UserId      uint      `gorm:"not null;column:user_id" json:"user_id"`
	ValidId     []byte    `gorm:"type:bytea;not null" json:"valid_id"`
	Selfie      []byte    `gorm:"type:bytea;not null" json:"selfie_id"`
	BackId      []byte    `gorm:"type:bytea;not null" json:"back_id"`
	Status      string    `gorm:"type:varchar(20);default:'pending'" json:"status"`
	SubmittedAt time.Time `gorm:"autoCreateTime" json:"submitted_at"`
	ReviewedAt  time.Time `json:"reviewed_at"`

	User User `gorm:"foreignKey:UserId;references:UserId" json:"user"`
}

// Updated UserNotification with TimeWithDate
type UserNotification struct {
	NotificationId int          `gorm:"primaryKey;column:notification_id" json:"notification_id"`
	Type           string       `gorm:"column:type" json:"type"`
	RequestId      int          `gorm:"column:request_id" json:"request_id"`
	FromUser       int          `gorm:"column:from_user" json:"from_user"`
	ToUser         int          `gorm:"column:to_user" json:"to_user"`
	Description    string       `gorm:"column:description" json:"description"`
	IsRead         bool         `gorm:"column:is_read" json:"is_read"`
	CreatedAt      TimeWithDate `gorm:"column:createdat" json:"createdat"`
}

// Updated UserNotification with TimeWithDate
type ChatNotification struct {
	NotificationId int          `gorm:"primaryKey;column:notification_id" json:"notification_id"`
	Type           string       `gorm:"column:type" json:"type"`
	FromUser       int          `gorm:"column:from_user" json:"from_user"`
	ToUser         int          `gorm:"column:to_user" json:"to_user"`
	Description    string       `gorm:"column:description" json:"description"`
	IsRead         bool         `gorm:"column:is_read" json:"is_read"`
	CreatedAt      TimeWithDate `gorm:"column:createdat" json:"createdat"`
}

// Claims structure remains the same
type Claims struct {
	UserId uint `json:"user_id"`
	jwt.StandardClaims
}

type ClientRepairmanConversation struct {
	ConversationId uint      `gorm:"primaryKey" json:"conversation_id"`
	ClientId       uint      `gorm:"not null" json:"client_id"`
	RepairmanId    uint      `gorm:"not null" json:"repairman_id"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"` // Add this line

	Client    User `gorm:"foreignKey:ClientId;references:UserId"`
	Repairman User `gorm:"foreignKey:RepairmanId;references:UserId"`
}

type ClientRepairmanMessage struct {
	MessageId      uint      `gorm:"primaryKey" json:"message_id"`
	ConversationId uint      `gorm:"not null" json:"conversation_id"`
	SenderId       uint      `gorm:"not null" json:"sender_id"`
	Message        string    `gorm:"type:text;not null" json:"message"`
	CreatedAt      time.Time `json:"created_at"`

	Conversation ClientRepairmanConversation `gorm:"foreignKey:ConversationId;references:ConversationId"`
	Sender       User                        `gorm:"foreignKey:SenderId;references:UserId"`
}

type GCashPayment struct {
	PaymentID     uint      `gorm:"primaryKey"`
	PaymentFrom   int       `gorm:"column:payment_from; not null"`
	PaymentTo     int       `gorm:"column:payment_to; not null"`
	TransactionId string    `gorm:"column:transaction_id; unique; not null"`
	Amount        float64   `gorm:"column:amount; not null"`
	GcashID       uint      `gorm:"column:gcash_id; not null"` // Foreign key to Gcash table
	PaymentDate   time.Time `gorm:"column:payment_date; not null"`
}

type Gcash struct {
	GcashID     uint   `gorm:"primaryKey"`
	UserID      uint   `gorm:"not null;unique"` // One GCash info per user
	GcashName   string `gorm:"type:varchar(100);not null"`
	GcashNumber string `gorm:"type:varchar(20);not null"`
	Createdat   time.Time
	Updatedat   time.Time
}

// TableName methods remain the same
func (User) TableName() string                   { return "users" }
func (Repairman) TableName() string              { return "users" }
func (EmailVer) TableName() string               { return "emailver" }
func (ServiceCategory) TableName() string        { return "service_categories" }
func (ServiceRequest) TableName() string         { return "service_requests" }
func (Review) TableName() string                 { return "reviews" }
func (UserVerification) TableName() string       { return "user_verifications" }
func (UserNotification) TableName() string       { return "user_notifications" }
func (ChatNotification) TableName() string       { return "chat_notifications" }
func (ClientRepairmanMessage) TableName() string { return "client_repairman_messages" }
func (GCashPayment) TableName() string           { return "gcash_payments" }
func (Gcash) TableName() string                  { return "gcash" }
