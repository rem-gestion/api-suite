package redis

import (
	"log"
	"sync"
	"users-api/src/config/envs"

	"github.com/go-redis/redis/v8"
)

var once sync.Once
var redisClient *redis.Client

func ConnectRedis() *redis.Client {
	REDIS_URI := envs.LoadEnvs(".env").Get("REDIS_URI")

	once.Do(func() {
		opt, err := redis.ParseURL(REDIS_URI)
		if err != nil {
			log.Printf("Error al parsear la URI de Redis: %v", err)
			return
		}

		redisClient = redis.NewClient(opt)

		if err := redisClient.Ping(redisClient.Context()).Err(); err != nil {
			log.Printf("Error al conectar con Redis: %v", err)
			redisClient = nil
			return
		}

		log.Println("Conexión a Redis establecida")
	})

	return redisClient
}
