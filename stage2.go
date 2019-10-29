package main

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/go-sql-driver/mysql"
	"github.com/satori/go.uuid"
	"io"
	"log"
	"net/http"
	"os"
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
)


var (
	mysqlClient *gorm.DB
	err error
)

func main() {

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

	router.GET("/", getIndexStg2)
	router.POST("/long2short", longToShortTransformStg2)
	router.POST("/short2long", shortToLongTransformStg2)

	router.Run()
}

func getIndexStg2(c *gin.Context) {
	pusher := c.Writer.Pusher()
	if pusher != nil {
		err := pusher.Push("/assets/app.jsx", nil)
		if err != nil {
			log.Printf("Failed to push: %v", err)
		}
	}
	c.HTML(http.StatusOK, "index.html", gin.H{})
}

func longToShortTransformStg2(c *gin.Context) {
	var urlPair UrlPair
	ret := INCORRECT_URL
	longUrl := c.DefaultPostForm(LONG_URL, DEFAULT_POST_FORM)

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
	c.JSON(SUCCESS_CODE, gin.H{
		SHORT_URL:    ret,
		STATUS: SUCCESS_CODE,
	})
}

func shortToLongTransformStg2(c *gin.Context) {
	var urlPair UrlPair
	var ret string
	shortUrl := c.DefaultPostForm(SHORT_URL, DEFAULT_POST_FORM)

	err := mysqlClient.Where(QUERY_SHORT_URL, shortUrl).First(&urlPair).Error
	if err != nil {
		ret = INCORRECT_URL
	} else {
		ret = urlPair.LongUrl
	}
	c.JSON(SUCCESS_CODE, gin.H{
		LONG_URL:    ret,
		STATUS: SUCCESS_CODE,
	})
}