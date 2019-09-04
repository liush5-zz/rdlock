package rdlock

import (
	"github.com/gomodule/redigo/redis"
	"log"
	"testing"
	"time"
)

var redisPool *redis.Pool

func init() {

	redisPool = &redis.Pool{
		MaxIdle:     10,
		MaxActive:   10,
		IdleTimeout: 240 * time.Second,
		Wait:        false,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", "0.0.0.0:6379",
				redis.DialPassword("redis"),
				redis.DialDatabase(10))
			if err != nil {
				return nil, err
			}
			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			if err != nil {
				log.Println("[ERROR] ping redis fail", err)
			}
			return err
		},
	}
}

func testTryLock(redisConn redis.Conn, key string, timeout int, sleepTime int) {
	lock, ok, err := TryLock(redisConn, key, timeout)
	if err != nil {
		log.Println("error", err)
		return
	}
	if !ok {
		log.Println("获取锁失败", lock)

	} else {
		log.Println("获取锁成功", lock)

	}
	defer lock.Unlock()

	time.Sleep(time.Second * time.Duration(sleepTime)) //sleepTime 秒后释放锁
}

func TestRedisLock(t *testing.T) {

	lockerKey := "test_redis_lock"
	timeOut := 5                                           //超时时间5s
	go testTryLock(redisPool.Get(), lockerKey, timeOut, 4) //4s后释放锁

	for i := 0; i < 8; i++ {
		testTryLock(redisPool.Get(), lockerKey, timeOut, 1)
	}

}
