package accounts

import (
	"errors"
	"strings"

	"gopkg.in/mgo.v2/bson"

	"github.com/thiepwong/microservices/common"
	"gopkg.in/mgo.v2"
)

type AccountRepository interface {
	Register(*RegisterModel) (interface{}, error)
	GetProfileById(sid uint64) (*Profile, error)
	GetUserById(username string) (*UserModel, error)

	CreateEmailPool(*EmailProfileModel) (bool, error)
	CreateMobilePool(*MobileProfileModel) (bool, error)
	UpdateProfile(*Profile) (bool, error)
}

type accountRepositoryContext struct {
	mgoSession *mgo.Session
	conf       *common.Config
}

func NewAccountReportsitory(db *mgo.Session, cfg *common.Config) AccountRepository {
	return &accountRepositoryContext{
		mgoSession: db,
		conf:       cfg,
	}
}

func (a *accountRepositoryContext) Register(data *RegisterModel) (interface{}, error) {
	err := a.mgoSession.DB(a.conf.Database.Mongo.Database).C("registers").Insert(data)
	if err != nil {
		if strings.Contains(err.Error(), "E11000") == true {
			return nil, errors.New("This username is registered before! Please use other username or re-activate this username")
		}
		return nil, err
	}
	return true, nil
}

func (a *accountRepositoryContext) GetProfileById(sid uint64) (*Profile, error) {
	var _profile Profile
	err := a.mgoSession.DB(a.conf.Database.Mongo.Database).C("profiles").FindId(sid).One(&_profile)
	if err != nil {
		return nil, errors.New("Profile not found!")
	}
	return &_profile, nil
}

func (a *accountRepositoryContext) GetUserById(username string) (*UserModel, error) {
	var _user UserModel
	_ = a.mgoSession.DB(a.conf.Database.Mongo.Database).C("users").FindId(username).One(&_user)

	return &_user, nil
}

func (a *accountRepositoryContext) CreateEmailPool(_em *EmailProfileModel) (bool, error) {
	_, err := a.mgoSession.DB(a.conf.Database.Mongo.Database).C("mailpools").Upsert(bson.M{"_id": _em.Email}, bson.M{"$set": bson.M{"code": _em.Code, "sid": _em.SID, "username": _em.Username, "email": _em.Email, "used": _em.Used, "full_name": _em.FullName}})
	if err != nil {
		return false, err
	}
	return true, nil
}

func (a *accountRepositoryContext) CreateMobilePool(_em *MobileProfileModel) (bool, error) {
	_, err := a.mgoSession.DB(a.conf.Database.Mongo.Database).C("mobilepools").Upsert(bson.M{"_id": _em.Mobile}, bson.M{"$set": bson.M{"code": _em.Code, "sid": _em.SID, "username": _em.Username, "mobile": _em.Mobile, "used": _em.Used, "full_name": _em.FullName}})
	if err != nil {
		return false, err
	}
	return true, nil
}

func (a *accountRepositoryContext) UpdateProfile(prof *Profile) (bool, error) {
	err := a.mgoSession.DB(a.conf.Database.Mongo.Database).C("profiles").UpdateId(prof.ID,
		bson.M{
			"$set": bson.M{
				"first_name": prof.FirstName,
				"last_name":  prof.LastName,
				"full_name":  prof.FullName,
				"gender":     prof.Gender,
				"address":    prof.Address,
				"birth_date": prof.BirthDate,
			},
		})
	if err != nil {
		return false, err
	}

	return true, nil
}
