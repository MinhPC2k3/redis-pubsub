package api

import (
	//"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	//"gorm.io/driver/postgres"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

type Product struct {
	gorm.Model
	Name  string `json:"name"`
	Price int    `json:"price"`
}

var GetData []Product

// var Initial = []Product{
// 	{Name: "car", Price: 100},
// 	{Name: "house", Price: 500},
// 	{Name: "computer", Price: 20},
// }
func InitRedisClient() *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	return rdb
}

//
func Request(c *gin.Context) {
	message := "received request"
	rdb := InitRedisClient()
	err := rdb.Publish(c, "myredis", "send me data").Err()
	if err != nil {
		panic(err)
	}
	c.IndentedJSON(http.StatusOK, message)

}
func CheckMessage() {
	var c *gin.Context
	rdb := InitRedisClient()
	sub := rdb.Subscribe(c, "myredis")
	for {
		message, err := sub.ReceiveMessage(c)
		if err != nil {
			panic(err)
		}

		fmt.Println(message)
	}
}
func ServeHTTP() {
	router := gin.Default()
	router.GET("/Request", Request)
	router.Run("localhost:8080")
}

// func OpenDB() *gorm.DB {
// 	dsn := "host=localhost user=gorm password=gorm dbname=gorm port=5432 sslmode=disable TimeZone=Asia/ShangHai"
// 	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
// 	if err != nil {
// 		panic(err)
// 	}
// 	db.AutoMigrate(&Product{})
// 	return db
// }
// db := OpenDB()
// var test Product
// db.Select("Name","Price").Create(&Initial)
// err := db.Table("products").Find(&GetData).Error
// if err != nil {
// 	fmt.Println(err)
// }
// fmt.Println(GetData)
// c.IndentedJSON(http.StatusOK, GetData)
