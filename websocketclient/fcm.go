package websocketclient

import (
	"context"
	"errors"
	"fixify_backend/model/users"
	"log"
	"sync"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/messaging"
	"gorm.io/gorm"
)

type FCMService struct {
	client *messaging.Client
	db     *gorm.DB
	mu     sync.Mutex
}

var (
	FCMInstance *FCMService
	initOnce    sync.Once
	initErr     error
)

// InitializeFCM should be called once at application startup
func InitializeFCM(db *gorm.DB, app *firebase.App) error {
	initOnce.Do(func() {
		client, err := app.Messaging(context.Background())
		if err != nil {
			initErr = err
			return
		}

		FCMInstance = &FCMService{
			client: client,
			db:     db,
		}
		log.Println("FCM service successfully initialized")
	})
	return initErr
}

func (f *FCMService) RegisterToken(userID uint, token string) error {
	if f == nil {
		return errors.New("FCM service not initialized")
	}

	f.mu.Lock()
	defer f.mu.Unlock()

	var user users.User
	if err := f.db.First(&user, userID).Error; err != nil {
		return err
	}

	user.FCMToken = token
	return f.db.Save(&user).Error
}

func SendPushNotification(toUserID uint, title, body string, data map[string]string) error {
	if FCMInstance == nil {
		log.Println("FCM Error: Service not initialized")
		return errors.New("FCM service not initialized")
	}

	FCMInstance.mu.Lock()
	defer FCMInstance.mu.Unlock()

	log.Printf("Looking up FCM token for user %d", toUserID)

	var user users.User
	if err := FCMInstance.db.First(&user, toUserID).Error; err != nil {
		log.Printf("FCM Error: User lookup failed for %d: %v", toUserID, err)
		return err
	}

	if user.FCMToken == "" {
		log.Printf("FCM Warning: No token for user %d (might need to register)", toUserID)
		return nil
	}

	log.Printf("Sending FCM to user %d with token: %s", toUserID, user.FCMToken)

	message := &messaging.Message{
		Token: user.FCMToken,
		Notification: &messaging.Notification{
			Title: title,
			Body:  body,
		},
		Data: data,
		Android: &messaging.AndroidConfig{
			Priority: "high", // Important for wakeup
		},
		APNS: &messaging.APNSConfig{
			Headers: map[string]string{
				"apns-priority": "10", // iOS high priority
			},
		},
	}

	response, err := FCMInstance.client.Send(context.Background(), message)
	if err != nil {
		log.Printf("FCM Error: Send failed to %d: %v", toUserID, err)
		return err
	}

	log.Printf("FCM Success: Message ID %s sent to %d", response, toUserID)
	return nil
}