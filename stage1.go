package main

import (
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"net/http"
	"os"
	"utils/URLShortener"
)

const (
	FIRST_ELEMENT = 0
	SUCCESS_CODE = 200
	INCORRECT_URL = "Incorrect url"
	LONG_URL = "longUrl"
	SHORT_URL = "shortUrl"
	STATUS = "status"
	DEFAULT_POST_FORM = "localhost"
	LOG_FILE_NAME = "MyURLShortenerLog.log"
)

var (
	long2shortMap map[string]string
)

func main() {

	f, _ := os.Create(LOG_FILE_NAME)
	gin.DefaultWriter = io.MultiWriter(f)

	long2shortMap = make(map[string]string)
	router := gin.Default()

	router.Static("/assets", "./assets")
	router.LoadHTMLFiles("templates/index.html")

	router.GET("/", getIndexStg1)
	router.POST("/long2short", longToShortTransformStg1)
	router.POST("/short2long", shortToLongTransformStg1)

	router.Run()
}

func getIndexStg1(c *gin.Context) {
	pusher := c.Writer.Pusher()
	if pusher != nil {
		err := pusher.Push("/assets/app.jsx", nil)
		if err != nil {
			log.Printf("Failed to push: %v", err)
		}
	}
	c.HTML(http.StatusOK, "index.html", gin.H{})
}

func longToShortTransformStg1(c *gin.Context) {
	longUrl := c.DefaultPostForm(LONG_URL, DEFAULT_POST_FORM)
	var ret string
	val, ok := long2shortMap[longUrl]
	if ok {
		ret = val
	} else {
		val, err := URLShortener.Transform(longUrl)
		if err != nil {
			ret = INCORRECT_URL
		} else {
			ret = val[FIRST_ELEMENT]
			long2shortMap[longUrl] = ret
		}
	}
	c.JSON(SUCCESS_CODE, gin.H{
		SHORT_URL:    ret,
		STATUS: SUCCESS_CODE,
	})
}

func shortToLongTransformStg1(c *gin.Context) {
	shortUrl := c.DefaultPostForm(SHORT_URL, DEFAULT_POST_FORM)
	ret := INCORRECT_URL
	for key, val := range long2shortMap {
		if val == shortUrl {
			ret = key
		}
	}
	c.JSON(SUCCESS_CODE, gin.H{
		LONG_URL:    ret,
		STATUS: SUCCESS_CODE,
	})
}