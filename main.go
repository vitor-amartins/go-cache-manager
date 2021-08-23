package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/joho/godotenv"
)

var ctx = context.Background()

type CacheManager struct {
	cache         *redis.Client
	id_expires_in time.Duration
}

func (cm *CacheManager) setIdByEmail(email string, id string) error {
	key := fmt.Sprintf("user_email:%s", email)
	return cm.cache.Set(ctx, key, id, cm.id_expires_in).Err()
}

func (cm *CacheManager) deleteIdByEmail(email string) error {
	key := fmt.Sprintf("user_email:%s", email)
	return cm.cache.Del(ctx, key).Err()
}

func (cm *CacheManager) getIdByEmail(email string) (string, error) {
	key := fmt.Sprintf("user_email:%s", email)
	return cm.cache.Get(ctx, key).Result()
}

func main() {
	godotenv.Load()

	host := os.Getenv("CACHE_HOST")
	port := os.Getenv("CACHE_PORT")
	addr := fmt.Sprintf("%s:%s", host, port)
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: "",
		DB:       0,
	})
	defer rdb.Close()

	cacheManager := CacheManager{
		cache:         rdb,
		id_expires_in: time.Hour * 24,
	}

	err := cacheManager.setIdByEmail("vitormartins@gravidadezero.space", "09b6701f-bc8d-41ba-8bde-66a03bfcbdfc")
	if err != nil {
		log.Panic(err)
	}

	id, err := cacheManager.getIdByEmail("vitormartins@gravidadezero.space")
	if err != nil {
		log.Panic(err)
	}
	fmt.Printf("%s\n", id)

	err = cacheManager.deleteIdByEmail("vitormartins@gravidadezero.space")
	if err != nil {
		log.Panic(err)
	}
}
