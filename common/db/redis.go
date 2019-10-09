package db

import (
	"fmt"
	"strconv"

	"github.com/thiepwong/microservices/common"

	"github.com/go-redis/redis"
)

//Redis struct for use
type Redis struct {
	Options redis.Options
	Client  *redis.Client
}

// RegisterRedis method
func registerRedis(red Redis) Redis {
	red.Client = redis.NewClient(&red.Options)
	return red
}

// GetRedisDb method
func GetRedisDb(cfg *common.CfgRd) *Redis {
	var _red Redis
	var e error
	_red.Options.Addr = fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
	_red.Options.Password = cfg.Password
	_red.Options.DB, e = strconv.Atoi(cfg.Database)
	if e != nil {
		_red.Options.DB = 0
	}
	_re := registerRedis(_red)
	fmt.Println(cfg.Host, cfg.Port, cfg.Password, cfg.Database)
	return &_re

}
