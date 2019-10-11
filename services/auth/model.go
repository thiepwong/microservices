package auth

import "time"

type Profile struct {
	ID        uint64
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
	ID             string          `json:"ID"  bson:"_id,omitempty" `
	Username       string          `json:"username" bson:"username"`
	Password       string          `json:"password" bson:"password"`
	SocialNetwork  []SocialNetwork `json:"SocialNetwork" bson:"social_network"`
	Profile        Profile         `json:"profile" bson:"profile"`
	RegisteredDate int64           `json:"registered_date" bson:"registered_date"`
	VerifyCode     string          `json:"verify_code" bson:"verify_code"`
	VerifiedDate   int             `json:"verified_date" bson:"verified_date"`
}

type SocialNetwork struct {
	Network string `json:"network" bson:"network"`
	Code    string `json:"code" bson:"code"`
}

type AccountModel struct {
	ID            uint64          `json:"ID"  bson:"_id,omitempty" `
	Username      string          `json:"username"  bson:"username"`
	Password      string          `json:"password" bson:"password"`
	SocialNetwork []SocialNetwork `json:"SocialNetwork" bson:"social_network"`
	ActivatedDate int64           `json:"activated_date" bson:"activated_date"`
	Profile       Profile         `json:"profile" bson:"profile"`
	Status        int             `json:"status" bson:"status"`
}

type SignInModel struct {
	Username string `json:"username" bson:"username"`
	Password string `json:"password" bson:"password"`
	System   string `json:"system" bson:"system"`
	Expired  int    `json:"expired"`
}

type ActivateModel struct {
	Username     string `json:"username"`
	ActivateCode string `json:"activate_code" bson:"activate_code"`
}

type SignInResponse struct {
	SmartID uint64 `json:"smart_id"`
	Token   string `json:"token"`
}

type UserModel struct {
	ID            string `json:"id"  bson:"_id,omitempty" `
	Username      string `json:"username"  bson:"username"`
	Password      string `json:"password" bson:"password"`
	ProfileID     uint64 `json:"profile_id" bson:"profile_id"`
	ActivatedDate int64  `json:"activated_date" bson:"activated_date"`
	Status        int    `json:"status" bson:"status"`
}

type SocialNetworkModel struct {
	ID            string `json:"id"  bson:"_id,omitempty" `
	Network       string `json:"network" bson:"network"`
	ProfileID     uint64 `json:"profile_id" bson:"profile_id"`
	ActivatedDate int64  `json:"activated_date" bson:"activated_date"`
	Status        int    `json:"status" bson:"status"`
}

type UserProfile struct {
	SmartID       uint64        `json:"smart_id"`
	Username      string        `json:"username"`
	Profile       *ProfileModel `json:"profile"`
	ActivatedDate int64         `json:"activated_date"`
	Status        int           `json:"status"`
}

type ProfileModel struct {
	ID        uint64 `json:"id"  bson:"_id,omitempty" `
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

type OtpModel struct {
	ID     string        `json:"id"`
	Mobile string        `json:"mobile"`
	TTL    time.Duration `json:"ttl"`
}

type UpdateContact struct {
	Contact string `json:"contact"`
	Code    string `json:"code"`
}
