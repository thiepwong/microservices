package auth

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/thiepwong/microservices/common/db"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const (
	salt = 11
)

type AuthRepository interface {
	SignIn(*SignInModel) (*UserModel, *ProfileModel, error)
	VerifyBySms(mobile string, otpCode string) (*RegisterModel, error)
	VerifyByEmail(email string, activateCode string) (*RegisterModel, error)
	CreateID(registerInfo *RegisterModel, smartID uint64) (*UserProfile, error)
}

type authRepositoryContext struct {
	db    *mgo.Database
	redis *db.Redis
}

func NewAuthRepository(db *mgo.Database, redis *db.Redis) AuthRepository {
	return &authRepositoryContext{
		db:    db,
		redis: redis,
	}
}

func (a *authRepositoryContext) SignIn(data *SignInModel) (*UserModel, *ProfileModel, error) {
	var _user UserModel
	var _profile ProfileModel

	err := a.db.C("users").FindId(data.Username).One(&_user)
	if err != nil {
		return nil, nil, errors.New("Not found this username in smart ID system")
	}
	err = a.db.C("profiles").FindId(_user.ProfileID).One(&_profile)
	if err != nil {
		return nil, nil, errors.New("Not found profile of this smart ID")
	}

	return &_user, &_profile, err
}

func (a *authRepositoryContext) VerifyByEmail(email string, activateCode string) (*RegisterModel, error) {

	var _register RegisterModel
	err := a.db.C("registers").FindId(email).One(&_register)
	if err != nil {
		return nil, errors.New("Username not found!")
	}

	if _register.VerifyCode == "" {
		return nil, errors.New("This account was activated before, please use forgoten password function!")
	}

	if _register.VerifyCode != activateCode {
		return nil, errors.New("Activate code is not match, please try again or resend the activate code!")
	}

	return &_register, nil
}

func (a *authRepositoryContext) VerifyBySms(mobile string, otpCode string) (register *RegisterModel, err error) {
	val, err := a.redis.Client.Get(otpCode).Result()
	if err != nil {
		return nil, err
	}
	var _otpObject OtpModel
	json.Unmarshal([]byte(val), &_otpObject)
	if _otpObject.Mobile != mobile {
		return nil, errors.New("OTP is invalid!")
	}

	err = a.db.C("registers").FindId(mobile).One(&register)
	if err != nil {
		return nil, errors.New("Username was not registered")
	}
	_ = a.redis.Client.Set(otpCode, _otpObject, 0).Err()
	return register, nil
}

func (a *authRepositoryContext) CreateID(registerInfo *RegisterModel, smartID uint64) (*UserProfile, error) {
	var _usr = &UserModel{}
	a.db.C("users").Find(bson.M{"username": registerInfo.Username}).One(_usr)

	if _usr.ProfileID > 0 {
		return nil, errors.New("This profile is activated before! Please re-verify profile if forgoten the password!")
	}

	_usr.ID = registerInfo.Username
	_usr.Username = registerInfo.Username
	_usr.Password = registerInfo.Password
	_usr.ProfileID = smartID
	_usr.ActivatedDate = time.Now().Unix()
	_usr.Status = 1

	var _profile = &ProfileModel{
		ID:        smartID,
		FirstName: registerInfo.Profile.FirstName,
		LastName:  registerInfo.Profile.LastName,
		FullName:  registerInfo.Profile.FullName,
		Gender:    registerInfo.Profile.Gender,
		Email:     registerInfo.Profile.Email,
		Mobile:    registerInfo.Profile.Mobile,
		BirthDate: registerInfo.Profile.BirthDate,
		Avatar:    registerInfo.Profile.Avatar,
		Address:   registerInfo.Profile.Address,
	}

	err := a.db.C("profiles").Insert(_profile)
	if err != nil {
		return nil, err
	}
	err = a.db.C("users").Insert(_usr)
	if err != nil {
		return nil, err
	}
	a.db.C("registers").Update(bson.M{"_id": _usr.Username}, bson.M{"$set": bson.M{"verified_date": _usr.ActivatedDate, "verify_code": nil}})

	_usr.Password = ""
	_sid := UserProfile{
		SmartID:       smartID,
		Username:      _usr.Username,
		Profile:       _profile,
		ActivatedDate: _usr.ActivatedDate,
		Status:        _usr.Status,
	}

	return &_sid, err
}
