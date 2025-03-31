package redis

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type RedisConfig struct {
	Host        string
	Password    string
	Port        int
	User        string
	DialTimeout time.Duration
	ReadTimeout time.Duration
}

const (
	envRedisHost = "REDIS_HOST"
	envRedisPort = "REDIS_PORT"
	envRedisUser = "REDIS_USER"
	envRedisPass = "REDIS_PASSWORD"
)

const servicePrefix = "awesome_service." // наш префикс сервиса

type RedisClient struct {
	config RedisConfig
	client *redis.Client
}

type Client interface {
	CheckJWTInBlacklist(ctx context.Context, jwtStr string) error
	WriteJWTToBlacklist(ctx context.Context, jwtStr string, jwtTTL time.Duration) error
}

func InitRedisConfig() RedisConfig {
	vp := viper.New()
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	port, err := strconv.Atoi(os.Getenv("REDIS_PORT"))
	if err != nil {
		return RedisConfig{}
	}

	return RedisConfig{
		Host:        os.Getenv("REDIS_HOST"),
		Password:    os.Getenv("REDIS_PASSWORD"),
		Port:        port,
		User:        os.Getenv(envRedisUser),
		DialTimeout: time.Duration(vp.GetInt("redis.dialTimeout")) * time.Second,
		ReadTimeout: time.Duration(vp.GetInt("redis.readTimeout")) * time.Second,
	}
}

func NewRedisClient(ctx context.Context, config RedisConfig) (*RedisClient, error) {
	client := &RedisClient{}

	client.config = config

	redisClient := redis.NewClient(&redis.Options{
		Password:    config.Password,
		Username:    config.User,
		Addr:        config.Host + ":" + strconv.Itoa(config.Port),
		DB:          0,
		DialTimeout: config.DialTimeout,
		ReadTimeout: config.ReadTimeout,
	})

	client.client = redisClient

	if _, err := redisClient.Ping(ctx).Result(); err != nil {
		return nil, fmt.Errorf("cant ping redis: %w", err)
	}

	return client, nil
}

func (c *RedisClient) Close() error {
	return c.client.Close()
}
