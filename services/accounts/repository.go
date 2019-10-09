package accounts

import (
	"errors"

	"gopkg.in/mgo.v2/bson"

	"gopkg.in/mgo.v2"
)

type AccountRepository interface {
	Register(*RegisterModel) (interface{}, error)
	AccountValidate(username string) (interface{}, error)
}

type accountRepositoryContext struct {
	db         *mgo.Database
	collection string
}

func NewAccountReportsitory(db *mgo.Database, coll string) AccountRepository {
	return &accountRepositoryContext{
		db:         db,
		collection: coll,
	}
}

func (a *accountRepositoryContext) Register(data *RegisterModel) (interface{}, error) {
	err := a.db.C("registers").Insert(data)
	if err != nil {
		return nil, err
	}
	return true, nil
}

func (a *accountRepositoryContext) AccountValidate(username string) (interface{}, error) {
	if username == "" {
		return nil, errors.New("Username is required!")
	}
	var _account RegisterModel
	err := a.db.C("accounts").Find(bson.M{"username": username}).One(&_account)
	return _account, err
}
