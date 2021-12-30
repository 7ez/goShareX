package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	LoadConfig()

	r := gin.Default()

	r.MaxMultipartMemory = 8 << 20 // 8 MB

	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  http.StatusOK,
			"message": "hello",
		})
	})
	r.GET("/i/:file", getFile)
	r.GET("/conf", genShareXConf)

	r.POST("/i/upload", uploadFile)

	r.Run(fmt.Sprintf(":%s", Config.Port))
}
