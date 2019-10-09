package accounts

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

// type SocialNetwork struct {
// 	Network string `json:"network" bson:"network"`
// 	Code    string `json:"code" bson:"code"`
// }

// type AccountModel struct {
// 	ID            int64   `json:"ID"  bson:"_id,omitempty" `
// 	Username      string  `json:"username"  bson:"username"`
// 	Password      string  `json:"password" bson:"password"`
// 	ActivatedDate int     `json:"activated_date" bson:"activated_date"`
// 	Profile       Profile `json:"profile" bson:"profile"`
// 	Status        string  `json:"status" bson:"status"`
// }
