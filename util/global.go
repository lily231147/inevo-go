package util

import (
	"context"

	"gorm.io/driver/mysql"

	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

var Ctx = context.Background()

// mysql 连接
var Db, _ = gorm.Open(
	mysql.Open("root:lily@tcp(localhost:3306)/jcyh?charset=utf8mb4&parseTime=True&loc=Local"),
	&gorm.Config{},
)

// Redis 连接
var Redis = redis.NewClient(&redis.Options{
	Addr:     "localhost:6379",
	Password: "", // no password set
	DB:       0,  // use default DB
})
