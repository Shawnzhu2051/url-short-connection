package main

import (
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"github.com/jinzhu/gorm"
	"github.com/satori/go.uuid"
	_ "github.com/go-sql-driver/mysql"
	"io"
	"log"
	"net/http"
	"os"
	"time"
	"utils/URLShortener"
)

type UrlPair struct {
	ID uuid.UUID
	LongUrl string `gorm:"unique"`
	ShortUrl string `gorm:"unique"`
}

const (
	FIRST_ELEMENT = 0
	SUCCESS_CODE = 200
	INCORRECT_URL = "Incorrect url"
	QUERY_LONG_URL = "long_url = ?"
	QUERY_SHORT_URL = "short_url = ?"
	LONG_URL = "longUrl"
	SHORT_URL = "shortUrl"
	STATUS = "status"
	DEFAULT_POST_FORM = "localhost"
	LOG_FILE_NAME = "MyURLShortenerLog.log"
	//Mysql arguments
	DB_TYPE = "mysql"
	DB_ARGS = "****:******@tcp(127.0.0.1:3306)/****?charset=utf8&parseTime=True&loc=Local"
	EXPIRE_TIME = 10 * time.Second
	//Redis arguments
	REDIS_ADDRESS = "localhost:6379"
	LONG_URL_REDIS_DB = 0
	SHORT_URL_REDIS_DB = 1
)

var (
	mysqlClient *gorm.DB
	err error
	redisLongURLClient *redis.Client
	redisShortURLClient *redis.Client
)


func main() {
	redisLongURLClient = redis.NewClient(&redis.Options{
		Addr:     REDIS_ADDRESS,
		Password: "",
		DB:       LONG_URL_REDIS_DB,
	})

	redisShortURLClient = redis.NewClient(&redis.Options{
		Addr:     REDIS_ADDRESS,
		Password: "",
		DB:       SHORT_URL_REDIS_DB,
	})

	mysqlClient, err = gorm.Open(DB_TYPE, DB_ARGS)
	if err != nil {
		log.Printf("Failed to connect database: %v", err)
	}
	defer mysqlClient.Close()
	mysqlClient.AutoMigrate(&UrlPair{})

	f, _ := os.Create(LOG_FILE_NAME)
	gin.DefaultWriter = io.MultiWriter(f)

	router := gin.Default()
	router.Static("/assets", "./assets")
	router.LoadHTMLFiles("templates/index.html")

	router.GET("/", getIndexStg3)
	router.POST("/long2short", longToShortTransformStg3)
	router.POST("/short2long", shortToLongTransformStg3)

	router.Run()
}

func getIndexStg3(c *gin.Context) {
	pusher := c.Writer.Pusher()
	if pusher != nil {
		err := pusher.Push("/assets/app.jsx", nil)
		if err != nil {
			log.Printf("Failed to push: %v", err)
		}
	}
	c.HTML(http.StatusOK, "index.html", gin.H{})
}

func longToShortTransformStg3(c *gin.Context) {
	ret := INCORRECT_URL
	longUrl := c.DefaultPostForm(LONG_URL, DEFAULT_POST_FORM)

	val, err := redisLongURLClient.Get(longUrl).Result()
	if err == redis.Nil {
		var urlPair UrlPair
		err := mysqlClient.Where(QUERY_LONG_URL, longUrl).First(&urlPair).Error
		if err != nil {
			val, err := URLShortener.Transform(longUrl)
			if err != nil {
				log.Printf("Failed to transform url: %v", err)
			} else {
				ret = val[FIRST_ELEMENT]
				id, _ := uuid.NewV4()
				urlPair = UrlPair{id ,longUrl, ret}
				mysqlClient.Create(&urlPair)
			}
		} else {
			ret = urlPair.ShortUrl
		}
		redisLongURLClient.Set(longUrl, ret, EXPIRE_TIME)
	} else if err != nil {
		log.Printf("redis error: %v", err)
		panic(err)
	} else {
		redisLongURLClient.Expire(longUrl, EXPIRE_TIME)
		ret = val
	}
	c.JSON(SUCCESS_CODE, gin.H{
		SHORT_URL:    ret,
		STATUS: SUCCESS_CODE,
	})
}

func shortToLongTransformStg3(c *gin.Context) {
	ret := INCORRECT_URL
	shortUrl := c.DefaultPostForm(SHORT_URL, DEFAULT_POST_FORM)

	val, err := redisShortURLClient.Get(shortUrl).Result()
	if err == redis.Nil {
		var urlPair UrlPair
		err := mysqlClient.Where(QUERY_SHORT_URL, shortUrl).First(&urlPair).Error
		if err != nil {
			// do nothing, remain void
		} else {
			ret = urlPair.LongUrl
		}
		redisShortURLClient.Set(shortUrl, ret, EXPIRE_TIME)
	} else if err != nil {
		log.Printf("redis error: %v", err)
		panic(err)
	} else {
		redisShortURLClient.Expire(shortUrl, EXPIRE_TIME)
		ret = val
	}
	c.JSON(SUCCESS_CODE, gin.H{
		LONG_URL:    ret,
		STATUS: SUCCESS_CODE,
	})
}