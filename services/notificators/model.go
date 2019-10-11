package notificators

import "time"

// OtpModel model
type OtpModel struct {
	ID     string        `json:"id"`
	Mobile string        `json:"mobile"`
	TTL    time.Duration `json:"ttl"`
}

type MailModel struct {
	Mail    string `json:"email"`
	Subject string `json:"subject"`
	Content string `json:"content"`
	Type    int    `json:"type"`
	Lang    string `json:"lang"`
}

type SmsModel struct {
	Mobile  string        `json:"mobile"`
	Content string        `json:"content"`
	Type    int           `json:"type"`
	Lang    string        `json:"lang"`
	TTL     time.Duration `json:"ttl"`
}

type Profile struct {
	FirstName string `json:"first_name" bson:"first_name"`
	LastName  string `json:"last_name" bson:"last_name"`
	FullName  string `json:"full_name" bson:"full_name"`
	Gender    int    `json:"gender" bson:"gender"`
	BirthDate int64  `json:"birth_date" bson:"birth_date"`
	Address   string `json:"address" bson:"address"`
	Avatar    string `json:"avatar" bson:"avatar"`
	Mobile    string `json:"mobile" bson:"mobile"`
	Email     string `json:"email" bson:"email"`
}

type RegisterModel struct {
	ID             string  `json:"ID"  bson:"_id,omitempty" `
	Username       string  `json:"username" bson:"username"`
	Password       string  `json:"password" bson:"password"`
	Profile        Profile `json:"profile" bson:"profile"`
	RegisteredDate int64   `json:"registered_date" bson:"registered_date"`
	VerifyCode     string  `json:"verify_code" bson:"verify_code"`
	VerifiedDate   int     `json:"verified_date" bson:"verified_date"`
}

type IrisSignInResponse struct {
	AccessToken string        `json:"access_token"`
	TokenType   string        `json:"token_type"`
	Expired     time.Duration `json:"expires_in"`
	//Error       string        `json:"error"`
}

type IrisSentResponse struct {
	Code    string `json:"Code"`
	Message string `json:"Message"`
	Data    bool   `json:"Data"`
}

type Message struct {
	SmsId       string `json:"SmsId"`
	PhoneNumber string `json:"PhoneNumber"`
	Content     string `json:"Content"`
	ContentType string `json:"ContentType"`
}

type SenderModel struct {
	Brandname   string    `json:"Brandname"`
	SendingList []Message `json:"SendingList"`
}

type EmailProfileModel struct {
	ID       string `json:"id" bson:"_id,omitempty"`
	SID      uint64 `json:"sid" bson:"sid"`
	Code     string `json:"code" bson:"code"`
	FullName string `json:"full_name" bson:"full_name"`
	Username string `json:"username" bson:"username"`
	Email    string `json:"email" bson:"email"`
	Used     bool   `json:"used" bson:"used"`
}
