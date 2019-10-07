package db

import (
	"fmt"

	"github.com/thiepwong/microservices/common"
	"gopkg.in/mgo.v2"
)

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
	db := session.DB(cfg.Database)

	return db, nil
}

func GetSession() *mgo.Session {
	dialInfo, err := mgo.ParseURL("mongodb://171.244.49.164:2688/ucenter")
	s, e := mgo.DialWithInfo(dialInfo)
	if e != nil {
		panic(err)
	}
	return s
}
