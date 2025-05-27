package controller

import (
	"fixify_backend/model/users"
	"time"

	"gorm.io/gorm"
)

// CreateUserNotification creates and saves a generic user notification.
// You can reuse this for any feature by changing the type and description.
func CreateUserNotification(
	db *gorm.DB,
	notificationType string,
	requestId int,
	fromUserId int,
	toUserId int,
	description string,
) error {
	notification := users.UserNotification{
		Type:        notificationType,
		RequestId:   requestId,
		FromUser:    fromUserId,
		ToUser:      toUserId,
		Description: description,
		IsRead:      false,
		CreatedAt:   users.TimeWithDate(time.Now()),
	}

	return db.Create(&notification).Error
}

func CreateChatNotification(
	db *gorm.DB,
	notificationType string,
	fromUserId int,
	toUserId int,
	description string,
) error {
	notification := users.ChatNotification{
		Type:        notificationType,
		FromUser:    fromUserId,
		ToUser:      toUserId,
		Description: description,
		IsRead:      false,
		CreatedAt:   users.TimeWithDate(time.Now()),
	}

	return db.Create(&notification).Error
}
