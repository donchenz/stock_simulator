package main


import (
    "log"
	"github.com/go-redis/redis"
)

const redisQName = "stock_q"


func createClient() *redis.Client {
    client := redis.NewClient(&redis.Options{
        Addr:     "localhost:6379",
        Password: "",
        DB:       0,
    })

    _, err := client.Ping().Result(); if err != nil {
        log.Fatal("redis error", err)
    }

    return client
}

var client = createClient()

func popMessage() (string, bool) {

	value, err := client.LPop(redisQName).Result() 
    if err != nil {
        return "", true
    }
    return value, false
}