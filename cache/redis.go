package cache
import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"time"
	redisClient "gopkg.in/redis.v5"
)
const (
	KeyNotFound  = "key not found"
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
	_, err := r.Client.Set(key, value, duration).Result()
	if err != nil {
		return err
	}
	return nil
}
// Del delete key.
func (r *Redis) Del(key string) error {
	_, err := r.Client.Del(key).Result()
	if err != nil {
		return err
	}
	return nil
}
// DelMany delete many keys.
func (r *Redis) DelMany(keys []string) error {
	_, err := r.Client.Del(keys...).Result()
	if err != nil {
		return err
	}
	return nil
}
// Get get key.
func (r *Redis) Get(key string) (string, error) {
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
	if result.Val() != "OK"{
		err := fmt.Errorf(ErrorClearCache, result.Err())
		return err
	}
	return nil
}
// Close is responsible for closing redis connection
func (r *Redis) Close() {
	err := r.Client.Close()
	if err != nil {
		log.Println(err)
	}
}