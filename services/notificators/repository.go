package notificators

import (
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
