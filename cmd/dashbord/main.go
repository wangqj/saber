package main

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/contrib/static"
	"saber/crl"
	"encoding/json"
)

func main() {
	r := gin.Default()
	// Default Gin router

	// Serve the frontend via views folder
	r.Use(static.Serve("/", static.LocalFile("cmd/dashbord/asserts", true)))

	r.GET("/slots", func(c *gin.Context) {
		r := crl.GetNode()
		n, _ := json.Marshal(r)
		c.JSON(200, string(n))
	})

	r.Run(":9999") // listen and serve on 0.0.0.0:8080
}
