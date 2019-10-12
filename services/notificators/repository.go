package notificators

import (
	"encoding/base64"
	"fmt"
	"time"

	"gopkg.in/mgo.v2"

	"github.com/thiepwong/microservices/common/db"
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
}

type notificatorRepositoryContext struct {
	redis *db.Redis
	db    *mgo.Database
}

//NewNotificatorRepository method
func NewNotificatorRepository(db *mgo.Database, redis *db.Redis) NotificatorRepository {
	return &notificatorRepositoryContext{
		redis: redis,
		db:    db,
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

	err = n.db.C("registers").FindId(email).One(&register)
	if err != nil {
		return nil, err
	}

	return register, nil

}

func (n *notificatorRepositoryContext) ReadRegisterByUser(username string) (register *RegisterModel, err error) {
	err = n.db.C("registers").FindId(username).One(&register)
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
	err := n.db.C("mailpools").FindId(email).One(&_mail)
	if err != nil {
		return nil, err
	}
	return &_mail, nil
}

func (n *notificatorRepositoryContext) ReadMobilePool(mobile string) (*MobileProfileModel, error) {
	var _mobile MobileProfileModel
	err := n.db.C("mobilepools").FindId(mobile).One(&_mobile)
	if err != nil {
		return nil, err
	}
	return &_mobile, nil
}
