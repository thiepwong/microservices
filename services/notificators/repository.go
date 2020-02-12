package notificators

import (
	"encoding/base64"
	"errors"
	"fmt"
	"strings"
	"time"

	"gopkg.in/mgo.v2/bson"

	uuid "github.com/satori/go.uuid"
	"github.com/thiepwong/microservices/common"
	"github.com/thiepwong/microservices/common/db"
	"gopkg.in/mgo.v2"
)

// NotificatorRepository method
type NotificatorRepository interface {
	SaveOTP(key string, value interface{}, ttl time.Duration) (bool, error)
	RemoveOTP(*OtpModel) (bool, error)
	ReadOTP(key string) (string, error)
	ReadMailActivatedCode(email string) (*RegisterModel, error)
	ReadRegisterByUser(username string) (register *RegisterModel, err error)
	ReadIrisToken(brandName string) string
	WriteIrisToken(brandName string, token string, ttl time.Duration) (bool, error)

	ReadMailPool(email string) (*EmailProfileModel, error)
	ReadMobilePool(mobile string) (*MobileProfileModel, error)

	ReadUserProfile(string) (*UserModel, *ProfileModel, error)
	ReadUserAccount(string) (*UserModel, error)
}

type notificatorRepositoryContext struct {
	redis      *db.Redis
	mgoSession *mgo.Session
	conf       *common.Config
}

//NewNotificatorRepository method
func NewNotificatorRepository(sess *mgo.Session, redis *db.Redis, cfg *common.Config) NotificatorRepository {
	return &notificatorRepositoryContext{
		redis:      redis,
		mgoSession: sess,
		conf:       cfg,
	}
}

func (n *notificatorRepositoryContext) SaveOTP(key string, value interface{}, ttl time.Duration) (bool, error) {
	err := n.redis.Client.Set(key, value, ttl*1000000000).Err()
	if err != nil {
		return false, err
	}

	return true, nil
}

func (n *notificatorRepositoryContext) ReadOTP(key string) (string, error) {
	val, err := n.redis.Client.Get(key).Result()
	if err != nil {
		return "", err
	}
	return val, nil
}

func (n *notificatorRepositoryContext) RemoveOTP(otp *OtpModel) (bool, error) {
	return false, nil
}

func (n *notificatorRepositoryContext) ReadMailActivatedCode(email string) (register *RegisterModel, err error) {

	_code := strings.Replace(uuid.Must(uuid.NewV4(), errors.New("error")).String(), "-", "", -1)
	if err != nil {
		return nil, errors.New("Cannot create activate code!")
	}

	err = n.mgoSession.DB(n.conf.Database.Mongo.Database).C("registers").UpdateId(email, bson.M{"$set": bson.M{"verify_code": _code}})

	err = n.mgoSession.DB(n.conf.Database.Mongo.Database).C("registers").FindId(email).One(&register)
	if err != nil {
		return nil, errors.New("This email was not registered!")
	}

	return register, nil

}

func (n *notificatorRepositoryContext) ReadRegisterByUser(username string) (*RegisterModel, error) {
	var register *RegisterModel
	err := n.mgoSession.DB(n.conf.Database.Mongo.Database).C("registers").FindId(username).One(&register)
	if err != nil {
		return nil, err
	}
	return register, nil
}

func (n *notificatorRepositoryContext) ReadIrisToken(brandName string) string {
	code := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:sms_token", brandName)))
	val, err := n.redis.Client.Get(code).Result()
	if err != nil {
		return ""
	}
	return val
}

func (n *notificatorRepositoryContext) WriteIrisToken(brandName string, token string, ttl time.Duration) (bool, error) {
	code := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:sms_token", brandName)))
	err := n.redis.Client.Set(code, token, ttl*1000000000).Err()
	if err != nil {
		return false, err
	}

	return true, nil
}

func (n *notificatorRepositoryContext) ReadMailPool(email string) (*EmailProfileModel, error) {
	var _mail EmailProfileModel
	err := n.mgoSession.DB(n.conf.Database.Mongo.Database).C("mailpools").FindId(email).One(&_mail)
	if err != nil {
		return nil, err
	}
	return &_mail, nil
}

func (n *notificatorRepositoryContext) ReadMobilePool(mobile string) (*MobileProfileModel, error) {
	var _mobile MobileProfileModel
	err := n.mgoSession.DB(n.conf.Database.Mongo.Database).C("mobilepools").FindId(mobile).One(&_mobile)
	if err != nil {
		return nil, err
	}
	return &_mobile, nil
}

func (n *notificatorRepositoryContext) ReadUserProfile(username string) (*UserModel, *ProfileModel, error) {
	var _usr UserModel
	var _prf ProfileModel
	err := n.mgoSession.DB(n.conf.Database.Mongo.Database).C("users").FindId(username).One(&_usr)
	if err != nil {
		return nil, nil, errors.New("This username is not exist!")
	}
	_ = n.mgoSession.DB(n.conf.Database.Mongo.Database).C("profiles").FindId(_usr.ProfileID).One(&_prf)

	return &_usr, &_prf, nil
}

func (n *notificatorRepositoryContext) ReadUserAccount(username string) (*UserModel, error) {
	var _usr UserModel
	err := n.mgoSession.DB(n.conf.Database.Mongo.Database).C("users").FindId(username).One(&_usr)
	if err != nil {
		return nil, errors.New("This username is not exist!")
	}
	return &_usr, nil
}
