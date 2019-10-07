package accounts

type Profile struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	FullName  string `json:"full_name"`
	Gender    int    `json:"gender"`
	BirthDate int64  `json:"birth_date"`
	Address   string `json:"address"`
	Avatar    string `json:"avatar"`
	Mobile    string `json:"mobile"`
	Email     string `json:"email"`
}

type RegisterModel struct {
	ID             string          `json:"ID"  bson:"_id,omitempty" `
	Username       string          `json:"username"`
	Password       string          `json:"password"`
	SocialNetwork  []SocialNetwork `json:"SocialNetwork"`
	Profile        Profile         `json:"profile"`
	RegisteredDate int64           `json:"registered_date"`
	VerifyCode     string          `json:"verify_code"`
	VerifiedDate   int             `json:"verified_date"`
}

type SocialNetwork struct {
	Network string
	Code    string
}

type AccountModel struct {
	ID            int64           `json:"ID"  bson:"_id,omitempty" `
	Username      string          `json:"username"`
	Password      string          `json:"password"`
	SocialNetwork []SocialNetwork `json:"SocialNetwork"`
	ActivatedDate int             `json:"activated_date"`
	Status        string          `json:"verify_code"`
}
