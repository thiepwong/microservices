package auth

import (
	"errors"
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const (
	salt = 11
)

type AuthRepository interface {
	SignIn(*SignInModel) (*AccountModel, error)
	VerifyBySms(mobile string, otpCode string) (*RegisterModel, error)
	VerifyByEmail(email string, activateCode string) (*RegisterModel, error)
	CreateID(registerInfo *AccountModel) (*AccountModel, error)
}

type authRepositoryContext struct {
	db *mgo.Database

	collection string
}

func NewAuthRepository(db *mgo.Database, coll string) AuthRepository {
	return &authRepositoryContext{
		db:         db,
		collection: coll,
	}
}

func (a *authRepositoryContext) SignIn(data *SignInModel) (*AccountModel, error) {
	var _account AccountModel
	err := a.db.C("accounts").Find(bson.M{"username": data.Username}).One(&_account)
	return &_account, err
}

func (a *authRepositoryContext) VerifyByEmail(email string, activateCode string) (*RegisterModel, error) {

	var _register RegisterModel
	err := a.db.C("registers").FindId(email).One(&_register)
	if err != nil {
		return nil, errors.New("Username not found!")
	}
	if _register.VerifyCode != activateCode {
		return nil, errors.New("Activate code is not match, please try again or resend the activate code!")
	}

	return &_register, nil
}

func (a *authRepositoryContext) VerifyBySms(mobile string, otpCode string) (*RegisterModel, error) {
	return nil, nil
}

func (a *authRepositoryContext) CreateID(account *AccountModel) (*AccountModel, error) {
	var _account = &AccountModel{}
	a.db.C("accounts").Find(bson.M{"username": account.Username}).One(_account)

	if _account.ID > 0 {
		return nil, errors.New("This profile is activated before! Please re-verify profile if forgoten the password!")
	}
	err := a.db.C("accounts").Insert(account)
	if err != nil {
		return nil, err
	}
	a.db.C("registers").Update(bson.M{"_id": account.Username}, bson.M{"$set": bson.M{"verified_date": time.Now().Unix(), "verify_code": nil}})
	return account, err
}
