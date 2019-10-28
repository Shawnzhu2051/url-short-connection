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
	ID uuid.UUID `json:"id"`
	LongUrl string `json:"longUrl"`
	ShortUrl string `json:"shortUrl"`
}

const (
	FIRST_ELEMENT = 0
	SUCCESS_CODE = 200
	INCORRECT_URL = "Incorrect url"
	QUERY_LONG_URL = "long_url = ?"
	QUERY_SHORT_URL = "short_url = ?"
)


var (
	db *gorm.DB
	err error
)

func main() {

	db, err = gorm.Open("mysql", "root:Mysqldemima1@tcp(127.0.0.1:3306)/gorm?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		log.Printf("Failed to connect database: %v", err)
	}
	defer db.Close()
	db.AutoMigrate(&UrlPair{})

	f, _ := os.Create("MyURLShortenerLog.log")
	gin.DefaultWriter = io.MultiWriter(f)

	router := gin.Default()
	router.Static("/assets", "./assets")
	router.LoadHTMLFiles("templates/index.html")

	router.GET("/", getIndexStg2)
	router.POST("/long2short", longToShortTransformStg2)
	router.POST("/short2long", shortToLongTransformStg2)

	router.Run() // listen and serve on 0.0.0.0:8080
}

func getIndexStg2(c *gin.Context) {
	if pusher := c.Writer.Pusher(); pusher != nil {
		if err := pusher.Push("/assets/app.jsx", nil); err != nil {
			log.Printf("Failed to push: %v", err)
		}
	}
	c.HTML(http.StatusOK, "index.html", gin.H{
		"title": "My URL Shortener",
	})
}

func longToShortTransformStg2(c *gin.Context) {
	var urlPair UrlPair;
	ret := INCORRECT_URL
	longUrl := c.DefaultPostForm("longUrl", "localhost")

	err := db.Where(QUERY_LONG_URL, longUrl).First(&urlPair).Error;
	if err != nil {
		val, err := URLShortener.Transform(longUrl);
		if err != nil {
			log.Printf("Failed to transform url: %v", err)
		} else {
			ret = val[FIRST_ELEMENT]
			u1, _ := uuid.NewV4()
			urlPair = UrlPair{u1 ,longUrl, ret}
			db.Create(&urlPair)
		}
	} else {
		ret = urlPair.ShortUrl
	}
	c.JSON(SUCCESS_CODE, gin.H{
		"shortUrl":    ret,
		"status": SUCCESS_CODE,
	})
}

func shortToLongTransformStg2(c *gin.Context) {
	var urlPair UrlPair
	var ret string
	shortUrl := c.DefaultPostForm("shortUrl", "localhost")

	err := db.Where(QUERY_SHORT_URL, shortUrl).First(&urlPair).Error;
	if err != nil {
		ret = INCORRECT_URL
	} else {
		ret = urlPair.LongUrl
	}
	c.JSON(SUCCESS_CODE, gin.H{
		"longUrl":    ret,
		"status": SUCCESS_CODE,
	})
}