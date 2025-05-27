package users

import "time"

type Admin struct {
	AdminId   int       `gorm:"primaryKey;column:admin_id" json:"admin_id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"createdat" gorm:"column:createdat;autoCreateTime"`
	UpdatedAt time.Time `json:"updatedat" gorm:"column:updatedat;autoUpdateTime"`
	FCMToken  string    `gorm:"size:255" json:"fcm_token"`
}

// Explicitly map to the correct table
func (Admin) TableName() string {
	return "admins"
}

type Conversation struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Admin1ID  uint      `json:"admin1_id"`
	Admin2ID  uint      `json:"admin2_id"`
	Admin1    Admin     `gorm:"foreignKey:Admin1ID;references:AdminId" json:"admin1"`
	Admin2    Admin     `gorm:"foreignKey:Admin2ID;references:AdminId" json:"admin2"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Message struct {
	ID             uint         `gorm:"primaryKey" json:"id"`
	ConversationID uint         `json:"conversation_id"`
	SenderID       uint         `json:"sender_id"`
	Sender         Admin        `gorm:"foreignKey:SenderID;references:AdminId" json:"sender"`
	Message        string       `json:"message"`
	IsRead         bool         `json:"is_read"`
	CreatedAt      time.Time    `json:"created_at"`
	Conversation   Conversation `gorm:"foreignKey:ConversationID" json:"conversation"`
}
