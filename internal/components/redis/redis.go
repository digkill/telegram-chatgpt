package redis

import (
	"context"
	"github.com/digkill/telegram-chatgpt/internal/config"
	"github.com/go-redis/redis/v8"
	log "github.com/sirupsen/logrus"
	"strconv"
	"time"
)

type RedisComponent struct {
	config *config.RedisConfig
	client *redis.Client
}

func (component *RedisComponent) Init() bool {

	component.client = redis.NewClient(&redis.Options{
		Addr:     component.config.Host + ":" + strconv.Itoa(component.config.Port),
		Password: component.config.Password,
		DB:       component.config.Db,
	})

	return true
}

func (component *RedisComponent) GetClient() *redis.Client {
	return component.client
}

func (component *RedisComponent) SetData(key string, value string, expiration time.Duration) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	return component.GetClient().Set(ctx, key, value, expiration).Err()
}

func (component *RedisComponent) GetData(key string) string {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	val, err := component.GetClient().Get(ctx, key).Result()
	if err != nil {
		log.Errorf("There is an error when make 'getData' Error: " + err.Error())
	}
	return val
}

func (component *RedisComponent) HasData(key string) bool {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	val, err := component.GetClient().Exists(ctx, key).Result()
	if err != nil {
		log.Errorf("There is an error when make 'hasData' Error: " + err.Error())
	}
	return val == int64(1)
}

func NewRedis(config *config.RedisConfig) *RedisComponent {
	redisComponent := &RedisComponent{
		config: config,
	}
	redisComponent.Init()
	return redisComponent
}
