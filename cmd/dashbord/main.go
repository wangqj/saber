package main

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/contrib/static"
	"fmt"
	"saber/crl"
	"encoding/json"
)

func main() {
	r := gin.Default()
	// Default Gin router

	// Serve the frontend via views folder
	r.Use(static.Serve("/", static.LocalFile("cmd/dashbord/asserts", true)))
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	r.GET("/slots", func(c *gin.Context) {
		r := crl.GetNode()
		n, _ := json.Marshal(r)
		fmt.Println(c.Param("zipcode"))
		c.JSON(200, string(n))
	})
	//r.LoadHTMLGlob("cmd/dashbord/asserts/index.html")
	//fmt.Println(r.Static())
	//r.GET("/dd/:name", func(c *gin.Context) {
	//	name := c.Param("name")
	//	fmt.Println("ddd",name)
	//	c.HTML(http.StatusOK, "index.html",gin.H{
	//
	//	})
	//	c.JSON(200, gin.H{
	//		"message": "pong",
	//		"name":name,
	//	})
	//})

	r.Run(":9999") // listen and serve on 0.0.0.0:8080
}
