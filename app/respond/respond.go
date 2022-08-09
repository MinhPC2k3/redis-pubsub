package respond

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"log"
	"os"
	"strconv"
	"time"

	// "encoding/csv"
	"fmt"
	// "log"
	// "os"

	"github.com/go-redis/redis/v8"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	// "github.com/fatih/structs"
)

type Product struct {
	gorm.Model
	Name  string `json:"name"`
	Price int    `json:"price"`
}

var ctx = context.Background()
var db = OpenDB()

func OpenDB() *gorm.DB {
	dsn := "host=localhost user=gorm password=gorm dbname=gorm port=5432 sslmode=disable TimeZone=Asia/ShangHai"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	db.AutoMigrate(&Product{})
	return db
}
func CreateClient() *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	return rdb
}
func CreateCsv(list []Product) {
	FirstRow := [][]string{
		{"", "name", "price"},
	}
	csvFile, err := os.Create("products.csv")

	if err != nil {
		log.Fatalf("failed creating file: %s", err)
	}
	writer := csv.NewWriter(csvFile)
	writer.WriteAll(FirstRow)
	for _, value := range list {
		var row []string
		row = append(row, strconv.FormatInt(int64(value.ID), 10))
		row = append(row, value.Name)
		row = append(row, strconv.FormatInt(int64(value.Price), 10))
		writer.Write(row)
	}
	writer.Flush()
	csvFile.Close()
	fmt.Println("success in create csv file")
}

func ProcessRequest() {
	rdb := CreateClient()
	sub := rdb.Subscribe(ctx, "myredis")
	for {
		time.Sleep(time.Second)
		message, err := sub.ReceiveMessage(ctx)
		if err != nil {
			panic(err)
		}
		if message.String() == "Message<myredis: send me data>" {
			var list []Product
			err := db.Raw("SELECT * FROM products").Scan(&list).Error
			if err != nil {
				panic(err)
			}
			//fmt.Println("Product:", list)
			JsonType, err := json.Marshal(list)
			if err != nil {
				panic(err)
			}
			value, err := rdb.Get(ctx, string(JsonType)).Result()
			if value != "exist" {
				fmt.Println(list)
				CreateCsv(list)
				rdb.Set(ctx, string(JsonType), "exist", 5*time.Second)
			} else {
				fmt.Println("csv file already exist")
			}

		}
	}
}

// func CreateCsv() {
// 	FirstRow := [][]string{
// 		{"", "name", "price"},
// 	}
// 	var List []Product
// 	err := db.Raw("SELECT * FROM products").Scan(&List).Error
// 	if err != nil {
// 		panic(err)
// 	}
// 	fmt.Println("Product:", List)
// 	csvFile, err := os.Create("products.csv")

// 	if err != nil {
// 		log.Fatalf("failed creating file: %s", err)
// 	}
// 	writer := csv.NewWriter(csvFile)
// 	writer.WriteAll(FirstRow)
// 	for _, value := range List {
// 		var row []string
// 		row = append(row, strconv.FormatInt(int64(value.ID), 10))
// 		row = append(row, value.Name)
// 		row = append(row, strconv.FormatInt(int64(value.Price), 10))
// 		writer.Write(row)
// 	}
// 	writer.Flush()
// 	csvFile.Close()
// }
