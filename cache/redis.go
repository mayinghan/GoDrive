package cache

import (
	"time"

	"github.com/garyburd/redigo/redis"
)

var (
	pool      *redis.Pool
	redisHost = "127.0.0.1:6379"
	redisPass = "123456"
)

func emailVerifPool() *redis.Pool {
	return &redis.Pool{
		MaxIdle:     50,
		MaxActive:   30,
		IdleTimeout: 300 * time.Second,
		Dial: func (redis.Conn. error) {
			// 1. open the connection
			// 2. auth
			// 3. select db 0
		}
	}
}
