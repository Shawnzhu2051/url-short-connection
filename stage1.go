package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
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
	long2shortMap = make(map[string]string)
	router := gin.Default()

	router.LoadHTMLFiles("templates/index.html")
	router.GET("/index", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{
			"title": "Main website",
		})
	})

	router.POST("/long2short",func(c *gin.Context){
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
		})
	})

	router.POST("/short2long",func(c *gin.Context){
		shortUrl := c.DefaultPostForm("shortUrl", "localhost")
		ret := INCORRECT_URL
		for key, val := range long2shortMap {
			if val == shortUrl {
				ret = key;
			}
		}
		c.JSON(SUCCESS_CODE, gin.H{
			"longUrl":    ret,
		})
	})
	router.Run() // listen and serve on 0.0.0.0:8080
}