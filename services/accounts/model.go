package accounts

type Profile struct {
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

type RegisterModel struct {
	ID             string  `json:"ID"  bson:"_id,omitempty" `
	Username       string  `json:"username" bson:"username"`
	Password       string  `json:"password" bson:"password"`
	Profile        Profile `json:"profile" bson:"profile"`
	RegisteredDate int64   `json:"registered_date" bson:"registered_date"`
	VerifyCode     string  `json:"verify_code" bson:"verify_code"`
	VerifiedDate   int     `json:"verified_date" bson:"verified_date"`
}

type UserModel struct {
	ID            string `json:"id"  bson:"_id,omitempty" `
	Username      string `json:"username"  bson:"username"`
	Password      string `json:"password" bson:"password"`
	ProfileID     uint64 `json:"profile_id" bson:"profile_id"`
	ActivatedDate int64  `json:"activated_date" bson:"activated_date"`
	Status        int    `json:"status" bson:"status"`
}

type AuthUpdate struct {
	SID      uint64 `json:"smart_id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Mobile   string `json:"mobile"`
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

type MobileProfileModel struct {
	ID       string `json:"id" bson:"_id,omitempty"`
	SID      uint64 `json:"sid" bson:"sid"`
	Code     string `json:"code" bson:"code"`
	FullName string `json:"full_name" bson:"full_name"`
	Username string `json:"username" bson:"username"`
	Mobile   string `json:"email" bson:"email"`
	Used     bool   `json:"used" bson:"used"`
}
