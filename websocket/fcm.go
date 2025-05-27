package websocket

import (
	"context"
	"fixify_backend/model/users"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/messaging"
	"google.golang.org/api/option"
	"gorm.io/gorm"
)

type FCMService struct {
	client *messaging.Client
	db     *gorm.DB
}

var FCMInstance *FCMService

func InitializeFCM(db *gorm.DB, credentialsFile string) error {
	opt := option.WithCredentialsFile(credentialsFile)
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		return err
	}

	client, err := app.Messaging(context.Background())
	if err != nil {
		return err
	}

	FCMInstance = &FCMService{
		client: client,
		db:     db,
	}

	return nil
}

func (f *FCMService) RegisterToken(userID uint, token string) error {
	// Update or create the FCM token in database
	var user users.Admin
	if err := f.db.First(&user, userID).Error; err != nil {
		return err
	}

	user.FCMToken = token
	return f.db.Save(&user).Error
}

func (f *FCMService) SendPushNotification(toUserID uint, title, body string, data map[string]string) error {
	var user users.Admin
	if err := f.db.First(&user, toUserID).Error; err != nil {
		return err
	}

	if user.FCMToken == "" {
		return nil // No token registered
	}

	message := &messaging.Message{
		Token: user.FCMToken,
		Notification: &messaging.Notification{
			Title: title,
			Body:  body,
		},
		Data: data,
	}

	_, err := f.client.Send(context.Background(), message)
	return err
}