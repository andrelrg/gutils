package cache

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"time"

	redisClient "gopkg.in/redis.v5"
)

const (
	KeyNotFound     = "key not found"
	ErrorClearCache = "fail to clean cache %s"
)

// Redis struct to manage redis.
type Redis struct {
	Client *redisClient.Client
	db     Config
}

// NewRedis is responsible for building a redis struct instance
func NewRedis(config Config) (*Redis, error) {
	red := Redis{db: config}
	err := red.Connect()
	if err != nil {
		return nil, err
	}
	return &red, nil
}

// Connect connects on redis database
func (r *Redis) Connect() error {
	db, _ := strconv.Atoi(r.db.GetDatabase())
	r.Client = redisClient.NewClient(&redisClient.Options{
		Addr:     r.db.GetHost() + ":" + strconv.Itoa(r.db.GetPort()),
		Password: r.db.GetPassword(),
		DB:       db,
	})
	_, err := r.Client.Ping().Result()
	if err != nil {
		return err
	}
	return nil
}

// Set set key.
func (r *Redis) Set(key, value string, duration time.Duration) error {
	key = buildServiceKey(key)

	_, err := r.Client.Set(key, value, duration).Result()
	if err != nil {
		return err
	}
	return nil
}

// MSet sets multiple key values
func (r *Redis) MSet(keys []string, values []interface{}, duration time.Duration) error {
	var ifaces []interface{}
	pipe := r.Client.TxPipeline()
	for i := range keys {
		key := keys[i]
		key = buildServiceKey(key)

		ifaces = append(ifaces, key, values[i])
		pipe.Expire(key, duration)
	}

	if err := r.Client.MSet(ifaces...).Err(); err != nil {
		return err
	}

	if _, err := pipe.Exec(); err != nil {
		return err
	}
	return nil
}

// Del delete key.
func (r *Redis) Del(key string) error {
	key = buildServiceKey(key)

	_, err := r.Client.Del(key).Result()
	if err != nil {
		return err
	}
	return nil
}

// DelMany delete many keys.
func (r *Redis) DelMany(keys []string) error {
	for i := range keys {
		keys[i] = buildServiceKey(keys[i])
	}

	_, err := r.Client.Del(keys...).Result()
	if err != nil {
		return err
	}
	return nil
}

// Get get key.
func (r *Redis) Get(key string) (string, error) {
	key = buildServiceKey(key)

	value, err := r.Client.Get(key).Result()
	if err == redisClient.Nil {
		return "", errors.New(KeyNotFound)
	} else if err != nil {
		return "", err
	}
	return value, nil
}

// Exist test if key exists.
func (r *Redis) Exist(key string) (bool, error) {
	key = buildServiceKey(key)

	_, err := r.Client.Get(key).Result()
	if err == redisClient.Nil {
		return false, nil
	} else if err != nil {
		return false, err
	}
	return true, nil
}

// FlushAll clears all keys in the cache
func (r *Redis) FlushAll() error {
	result := r.Client.FlushAll()
	if result.Val() != "OK" {
		err := fmt.Errorf(ErrorClearCache, result.Err())
		return err
	}
	return nil
}

// Scan gets redis keys based on a match pattern, something like r.Scan("my_key:*")
func (r *Redis) Scan(match string) ([]string, error) {
	var err error
	var cursor uint64
	var keys []string
	var result []string
	var count int64

	count = 100

	for {
		keys, cursor, err = r.Client.Scan(cursor, match, count).Result()

		if err != nil {
			log.Println(fmt.Printf("error scanning cache keys by match %s, %s", match, err.Error()))
			return result, err
		}

		result = append(result, keys...)
		if cursor == 0 {
			break
		}
	}

	return result, err
}

// DelByPattern gets redis keys based on a match pattern using Scan() method and then,
// using DelMany(), removes these keys from redis cache
func (r *Redis) DelByPattern(match string) error {
	var err error
	var keys []string

	keys, err = r.Scan(match)

	if len(keys) == 0 {
		return err
	}

	err = r.DelMany(keys)
	return err
}

// DelByKeysPattern DO NOT USE THIS METHOD unless you're absolutely sure about what you're doing!
// it gets redis keys based on a match pattern using Keys() method and then,
// using DelMany(), removes these keys from redis cache
func (r *Redis) DelByKeysPattern(match string) error {
	var err error
	var keys []string

	log.Println("please, DON'T USE DelByKeysPattern() method on production environment unless you're 100% sure")

	keys, err = r.Client.Keys(match).Result()

	if err != nil {
		return err
	}

	err = r.DelMany(keys)
	return err
}

// Close is responsible for closing redis connection
func (r *Redis) Close() {
	err := r.Client.Close()
	if err != nil {
		log.Println(err)
	}
}

//build key from binary name
func buildServiceKey(key string) string {
	binaryName := os.Args[0]
	binaryName = filepath.Base(binaryName)

	return binaryName + ":" + key
}
