package notificators

type NotificatorService interface {
	SendEmail(targetMail string, subject string, content string) (bool, error)
	SendSMS(mobile string, content string) (interface{}, error)
	SendFirebase(channel string, title string, content string) (interface{}, error)
}
