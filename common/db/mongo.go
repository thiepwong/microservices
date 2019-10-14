package db

import (
	"fmt"

	"github.com/thiepwong/microservices/common"
	"gopkg.in/mgo.v2"
)

type MgoHelper struct {
	Addr      []string
	User      string
	Password  string
	PoolLimit int
	FailFast  bool
	Session   *mgo.Session
	Source    string
}

// Connect connect to mongo server
func (mh *MgoHelper) Connect(cfg *common.CfgMg) error {
	di := &mgo.DialInfo{
		Addrs:     mh.Addr,
		Username:  mh.User,
		Password:  mh.Password,
		PoolLimit: mh.PoolLimit,
		FailFast:  mh.FailFast,
	}

	s, err := mgo.DialWithInfo(di)
	if err != nil {
		return err
	}
	defer s.Close()

	mh.Session = s.Copy()

	return nil
}

func GetMongoDb(cfg *common.CfgMg) (*mgo.Database, error) {

	connectAdd := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
	dialInfo1 := &mgo.DialInfo{
		Addrs:    []string{connectAdd},
		Username: cfg.Username,
		Password: cfg.Password,
		Database: cfg.Database,
		Source:   cfg.Auth,
	}

	// session, err := mgo.Dial(host)

	//dialInfo, err := mgo.ParseURL("mongodb://171.244.49.164:2688/ucenter")

	session, err := mgo.DialWithInfo(dialInfo1)

	if err != nil {
		return nil, err
	}

	defer session.Close()
	db := session.DB(cfg.Database)
	return db, nil
}

func GetSession(cfg *common.CfgMg) *mgo.Session {
	connectAdd := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
	dialInfo := &mgo.DialInfo{
		Addrs:    []string{connectAdd},
		Username: cfg.Username,
		Password: cfg.Password,
		Database: cfg.Database,
		Source:   cfg.Auth,
	}
	s, e := mgo.DialWithInfo(dialInfo)
	if e != nil {
		panic(e)
	}
	defer s.Close()
	return s.Copy()
}
