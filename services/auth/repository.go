package auth

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/thiepwong/microservices/common"
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
	ReadMailPool(email string, code string) (*EmailProfileModel, error)
	ReadMobilePool(mobile string, code string) (*MobileProfileModel, error)
	UpdateProfileWithCombineUser(username string, smartID uint64, contactType int) (bool, error)
	UpdateProfile(contact string, smartID uint64, contactType int) (bool, error)
}

type authRepositoryContext struct {
	mgoSession *mgo.Session
	redis      *db.Redis
	conf       *common.Config
}

func NewAuthRepository(sess *mgo.Session, redis *db.Redis, cfg *common.Config) AuthRepository {
	return &authRepositoryContext{
		mgoSession: sess,
		redis:      redis,
		conf:       cfg,
	}
}

func (a *authRepositoryContext) SignIn(data *SignInModel) (*UserModel, *ProfileModel, error) {
	var _user UserModel
	var _profile ProfileModel

	err := a.mgoSession.DB(a.conf.Database.Mongo.Database).C("users").FindId(data.Username).One(&_user)
	if err != nil {
		return nil, nil, errors.New("Not found this username in smart ID system")
	}
	err = a.mgoSession.DB(a.conf.Database.Mongo.Database).C("profiles").FindId(_user.ProfileID).One(&_profile)
	if err != nil {
		return nil, nil, errors.New("Not found profile of this smart ID")
	}

	return &_user, &_profile, err
}

func (a *authRepositoryContext) VerifyByEmail(email string, activateCode string) (*RegisterModel, error) {

	var _register RegisterModel
	err := a.mgoSession.DB(a.conf.Database.Mongo.Database).C("registers").FindId(email).One(&_register)
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

	err = a.mgoSession.DB(a.conf.Database.Mongo.Database).C("registers").FindId(mobile).One(&register)
	if err != nil {
		return nil, errors.New("Username was not registered")
	}
	_ = a.redis.Client.Set(otpCode, _otpObject, 0).Err()
	return register, nil
}

func (a *authRepositoryContext) CreateID(registerInfo *RegisterModel, smartID uint64) (*UserProfile, error) {
	var _usr = &UserModel{}
	a.mgoSession.DB(a.conf.Database.Mongo.Database).C("users").Find(bson.M{"username": registerInfo.Username}).One(_usr)

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

	err := a.mgoSession.DB(a.conf.Database.Mongo.Database).C("profiles").Insert(_profile)
	if err != nil {
		return nil, err
	}
	err = a.mgoSession.DB(a.conf.Database.Mongo.Database).C("users").Insert(_usr)
	if err != nil {
		return nil, err
	}
	a.mgoSession.DB(a.conf.Database.Mongo.Database).C("registers").Update(bson.M{"_id": _usr.Username}, bson.M{"$set": bson.M{"verified_date": _usr.ActivatedDate, "verify_code": nil}})

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

func (a *authRepositoryContext) ReadMailPool(email string, code string) (*EmailProfileModel, error) {

	var _emP EmailProfileModel
	err := a.mgoSession.DB(a.conf.Database.Mongo.Database).C("mailpools").FindId(email).One(&_emP)
	if err != nil {
		return nil, errors.New("Email not found in pool, please re-add email to profile")
	}

	return &_emP, nil
}

func (a *authRepositoryContext) ReadMobilePool(mobile string, code string) (*MobileProfileModel, error) {
	var _mb MobileProfileModel
	err := a.mgoSession.DB(a.conf.Database.Mongo.Database).C("mobilepools").FindId(mobile).One(&_mb)
	if err != nil {
		return nil, errors.New("Mobile not found in pool, please re-add mobile to profile")
	}

	return &_mb, nil
}

func (a *authRepositoryContext) UpdateProfileWithCombineUser(username string, smartID uint64, contactType int) (bool, error) {

	var err error
	err = a.mgoSession.DB(a.conf.Database.Mongo.Database).C("users").Update(bson.M{"_id": username}, bson.M{"$set": bson.M{"profile_id": smartID}})
	if err != nil {
		return false, err
	}

	switch contactType {
	case 1:
		//Update with combine user as Email
		err = a.mgoSession.DB(a.conf.Database.Mongo.Database).C("profiles").Update(bson.M{"_id": smartID}, bson.M{"$set": bson.M{"email": username}})

		break
	case 2:
		// Update with combine user as Mobile
		err = a.mgoSession.DB(a.conf.Database.Mongo.Database).C("profiles").Update(bson.M{"_id": smartID}, bson.M{"$set": bson.M{"mobile": username}})
		break
	}

	if err != nil {
		return false, err
	}

	return true, nil
}

func (a *authRepositoryContext) UpdateProfile(contact string, smartID uint64, contactType int) (bool, error) {
	var err error

	switch contactType {
	case 1:
		//Update with combine user as Email
		err = a.mgoSession.DB(a.conf.Database.Mongo.Database).C("profiles").Update(bson.M{"_id": smartID}, bson.M{"$set": bson.M{"email": contact}})

		break
	case 2:
		// Update with combine user as Mobile
		err = a.mgoSession.DB(a.conf.Database.Mongo.Database).C("profiles").Update(bson.M{"_id": smartID}, bson.M{"$set": bson.M{"mobile": contact}})
		break
	}

	if err != nil {
		return false, err
	}

	return false, nil
}
