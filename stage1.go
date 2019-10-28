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
)

var (
	long2shortMap map[string]string
)

func main() {

	f, _ := os.Create("MyURLShortenerLog.log")
	gin.DefaultWriter = io.MultiWriter(f)

	long2shortMap = make(map[string]string)
	router := gin.Default()

	router.Static("/assets", "./assets")
	router.LoadHTMLFiles("templates/index.html")

	router.GET("/", getIndexStg1)
	router.POST("/long2short", longToShortTransformStg1)
	router.POST("/short2long", shortToLongTransformStg1)

	router.Run() // listen and serve on 0.0.0.0:8080
}

func getIndexStg1(c *gin.Context) {
	if pusher := c.Writer.Pusher(); pusher != nil {
		if err := pusher.Push("/assets/app.jsx", nil); err != nil {
			log.Printf("Failed to push: %v", err)
		}
	}
	c.HTML(http.StatusOK, "index.html", gin.H{
		"title": "Main website",
	})
}

func longToShortTransformStg1(c *gin.Context) {
	longUrl := c.DefaultPostForm("longUrl", "localhost")
	var ret string
	if val, ok := long2shortMap[longUrl]; ok {
		ret = val
	} else {
		val, err := URLShortener.Transform(longUrl);
		if err != nil {
			ret = INCORRECT_URL
		} else {
			ret = val[FIRST_ELEMENT]
			long2shortMap[longUrl] = ret
		}
	}
	c.JSON(SUCCESS_CODE, gin.H{
		"shortUrl":    ret,
		"status": SUCCESS_CODE,
	})
}

func shortToLongTransformStg1(c *gin.Context) {
	shortUrl := c.DefaultPostForm("shortUrl", "localhost")
	ret := INCORRECT_URL
	for key, val := range long2shortMap {
		if val == shortUrl {
			ret = key;
		}
	}
	c.JSON(SUCCESS_CODE, gin.H{
		"longUrl":    ret,
		"status": SUCCESS_CODE,
	})
}