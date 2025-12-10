package websocket

const (
	NotificationTypeNewDuck = "new_duck_created"
)

type Notification struct {
	Type string      `json:"type"`
	Data interface{} `json:"data"`
}

func NewNotification(notificationType string, data interface{}) *Notification {
	return &Notification{
		Type: notificationType,
		Data: data,
	}
}
