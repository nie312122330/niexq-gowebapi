package redisext

import (
	"errors"
	"fmt"
	"niexq-gowebapi/ginext"

	"github.com/gomodule/redigo/redis"
	"github.com/nie312122330/niexq-gotools/logext"
	"go.uber.org/zap"
)

var logger *zap.Logger

// RedisService ...
type RedisService struct {
	RedisPool *redis.Pool
}

func init() {
	logger = logext.DefaultLogger("redis")
}

// PutStr ...
func (service *RedisService) PutStr(key string, val string) error {
	conn := service.RedisPool.Get()
	defer conn.Close()
	resp, err := redis.String(conn.Do("SET", key, val))
	if nil != err {
		return err
	}
	if "OK" != resp {
		return errors.New("未返回OK")
	}
	return nil
}

// PutExStr ...
func (service *RedisService) PutExStr(key string, val string, sencond int) error {
	conn := service.RedisPool.Get()
	defer conn.Close()
	resp, err := redis.String(conn.Do("SETEX", key, sencond, val))
	if nil != err {
		return err
	}
	if "OK" != resp {
		return errors.New("未返回OK")
	}
	return nil
}

// PutNxExStr ...
func (service *RedisService) PutNxExStr(key string, val string, sencond int) error {
	conn := service.RedisPool.Get()
	defer conn.Close()
	resp, err := redis.String(conn.Do("SET", key, val, "EX", sencond, "NX"))
	if nil != err {
		return err
	}
	if "OK" != resp {
		return errors.New("未返回OK")
	}
	return nil
}

// GetStr ...
func (service *RedisService) GetStr(key string) (string, error) {
	conn := service.RedisPool.Get()
	defer conn.Close()
	val, err := redis.String(conn.Do("GET", key))
	if nil != err {
		return "", err
	}
	return val, nil
}

// ExpireKey ...
func (service *RedisService) ExpireKey(key string, sencond int) error {
	conn := service.RedisPool.Get()
	defer conn.Close()
	resp, err := redis.Int64(conn.Do("EXPIRE", key, sencond))
	if nil != err {
		return err
	}
	if resp <= 0 {
		err := ginext.RunTimeError{Err: "设置成功数小于0"}
		return &err
	}
	return nil
}

// ClearByKeyPrefix 清理指定前缀的KEY
func (service *RedisService) ClearByKeyPrefix(keyPrefix string) (int, error) {
	conn := service.RedisPool.Get()
	defer conn.Close()
	keyPattner := fmt.Sprintf("%s*", keyPrefix)
	//扫描Key
	keys, err := scanKeysWithConn(conn, 0, keyPattner, nil, 1000)
	if nil != err {
		return 0, err
	}
	//删除
	if len(keys) > 0 {
		var delKeys = make([]interface{}, len(keys))
		for i := 0; i < len(keys); i++ {
			delKeys[i] = keys[i]
		}
		return redis.Int(conn.Do("DEL", delKeys...))
	}
	return 0, nil
}

func scanKeysWithConn(conn redis.Conn, cur int, keyPattner string, lastKeys []string, maxLen int) ([]string, error) {
	reply, err := conn.Do("SCAN", cur, "MATCH", keyPattner, "COUNT", maxLen)
	if nil == err {
		replyArray := reply.([]interface{})
		cur, _ = redis.Int(replyArray[0], nil)
		curKeys, _ := redis.Strings(replyArray[1], nil)
		var keys []string
		if nil != lastKeys {
			keys = append(lastKeys, curKeys...)
		} else {
			keys = curKeys
		}
		if len(keys) > maxLen {
			return nil, &ginext.RunTimeError{Err: fmt.Sprintf("Key数量超过了%d", maxLen)}
		}
		if cur != 0 {
			return scanKeysWithConn(conn, cur, keyPattner, keys, maxLen)
		}
		return keys, nil
	}
	return nil, err
}
