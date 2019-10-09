package notificators

import "time"

// OtpModel model
type OtpModel struct {
	ID     string        `json:"id"`
	Mobile string        `json:"mobile"`
	TTL    time.Duration `json:"ttl"`
}

type RegisterModel struct {
	ID             string `json:"ID"  bson:"_id,omitempty" `
	Username       string `json:"username" bson:"username"`
	RegisteredDate int64  `json:"registered_date" bson:"registered_date"`
	VerifyCode     string `json:"verify_code" bson:"verify_code"`
}

type MailModel struct {
	Mail    string `json:"email"`
	Subject string `json:"subject"`
	Content string `json:"content"`
	Type    int    `json:"type"`
}
