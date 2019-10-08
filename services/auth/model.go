package auth

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
	Expired  int64  `json:"expired"`
}

type ActivateModel struct {
	Username     string `json:"username"`
	ActivateCode string `json:"activate_code" bson:"activate_code"`
}

type SignInResponse struct {
	SmartID uint64 `json:"smart_id"`
	Token   string `json:"token"`
}
