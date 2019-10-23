package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"utils/URLShortener"
)

var (
	long2shortMap map[string]string
)

func main() {
	long2shortMap = make(map[string]string)
	router := gin.Default()

	router.POST("/long2short",func(c *gin.Context){
		longUrl := c.DefaultPostForm("longUrl", "localhost")
		var ret string
		if val, ok := long2shortMap[longUrl]; ok {
			ret = val
		} else {
			val, err := URLShortener.Transform(longUrl);
			if err != nil {
				ret = "Incorrect url"
			} else {
				ret = val[0]
				fmt.Println(longUrl)
				fmt.Println(ret)
				long2shortMap[longUrl] = ret
			}
		}
		c.JSON(200, gin.H{
			"status":  "posted",
			"shortUrl":    ret,
		})
	})

	router.POST("/short2long",func(c *gin.Context){
		shortUrl := c.DefaultPostForm("shortUrl", "localhost")
		ret := "Unrecognized url"
		for key, val := range long2shortMap {
			if val == shortUrl {
				ret = key;
			}
		}
		c.JSON(200, gin.H{
			"status":  "posted",
			"longUrl":    ret,
		})
	})
	router.Run(":8080") // listen and serve on 0.0.0.0:8080
}