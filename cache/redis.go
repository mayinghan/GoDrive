package cache

import (
	"fmt"
	"time"

	"github.com/garyburd/redigo/redis"
)

var (
	emailVeriPool *redis.Pool
	redisHost     = "127.0.0.1:6379"
	redisPass     = "123456"
	emailVeriDB   = 0
)

func createPool(dbtype int) *redis.Pool {
	return &redis.Pool{
		MaxIdle:     50,
		MaxActive:   30,
		IdleTimeout: 300 * time.Second,
		Dial: func() (redis.Conn, error) {
			// 1. open the connection
			c, err := redis.Dial("tcp", redisHost)
			if err != nil {
				panic(err)
			}

			// 2. auth
			if _, err = c.Do("AUTH", redisPass); err != nil {
				c.Close()
				fmt.Println("Redis auth failed")
				return nil, err
			}
			// 3. select db 0
			if _, err = c.Do("SELECT", dbtype); err != nil {
				c.Close()
				fmt.Println("Select redis db for emailveri failed")
				panic(err)
			}
			return c, nil
		},
		TestOnBorrow: func(conn redis.Conn, t time.Time) error {
			// if the time cost since connection is less than 5 minutes
			if time.Since(t) < 5*time.Minute {
				return nil
			}
			_, err := conn.Do("PING")
			return err
		},
	}
}

func init() {
	emailVeriPool = createPool(emailVeriDB)
}

// EmailVeriPool returns the redis pool for email verification
func EmailVeriPool() *redis.Pool {
	return emailVeriPool
}
