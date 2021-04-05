package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

type access struct {
	barrierIn  string
	barrierOut string
}

func main() {
	r := gin.Default()
	a := make(map[string]*access)

	r.POST("/barrier/:barrierID/in/:registration", func(c *gin.Context) {
		reg := c.Param("registration")
		bID := c.Param("barrierID")

		a[reg] = &access{barrierIn: bID}

		c.JSON(200, gin.H{
			"message": fmt.Sprintf("registration-id: %s entered the parking lot using barrier %s", reg, bID),
		})
	})

	r.POST("/barrier/:barrierID/out/:registration", func(c *gin.Context) {
		reg := c.Param("registration")
		bID := c.Param("barrierID")

		car := a[reg]
		car.barrierOut = bID
		delete(a, reg)

		c.JSON(200, gin.H{
			"message": fmt.Sprintf("registration-id: %s exited the parking lot using barrier %s", reg, bID),
		})
	})

	r.GET("/overview", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": fmt.Sprintf("there are %v cars in the parking lot", len(a)),
		})
	})

	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
