package main

import (
	"./hand"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {

	r := gin.Default()
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:3000", "http://0.0.0.0:12312"}
	config.AllowCredentials = true
	config.AllowMethods = []string{"GET", "POST"}
	r.Use(cors.New(config))
	r.LoadHTMLFiles("index.html")

	r.GET("/", func(c *gin.Context) {
		c.HTML(200, "index.html", nil)
	})

	r.Static("/js", ".")

	r.GET("/socket.io/", gin.WrapH(hand.Server()))

	r.Run("localhost:12312")
}
