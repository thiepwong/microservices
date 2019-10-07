package auth

import (
	"errors"

	"github.com/thiepwong/microservices/common"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type AuthRepository interface {
	SignIn(*SignInModel) (interface{}, error)
}

type authRepositoryContext struct {
	db         *mgo.Database
	collection string
}

func NewAuthRepository(db *mgo.Database, coll string) AuthRepository {
	return &authRepositoryContext{
		db:         db,
		collection: coll,
	}
}

func (a *authRepositoryContext) SignIn(data *SignInModel) (interface{}, error) {
	var _account AccountModel
	err := a.db.C("register").Find(bson.M{"username": data.Username}).One(&_account)
	if err != nil {
		return nil, err
	}
	valid := common.PasswordCompare(data.Password, _account.Password, 11)
	if valid == false {
		return nil, errors.New("Password is invalid!")

	} else {
		return _account, nil
	}
}
