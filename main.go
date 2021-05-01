package main

import (
	"fmt"
	"log"

	"github.com/RicardoLinck/parking-lot-api/barrier"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	bc := barrier.NewBarrierConfig("./logs/barriers")

	r.POST("/barrier/:barrierID/in/:registration", func(c *gin.Context) {
		reg := c.Param("registration")
		bID := c.Param("barrierID")
		if err := bc.Validate(bID); err != nil {
			c.JSON(404, gin.H{
				"error": err,
			})
		}

		err := bc.In(bID, reg)
		if err != nil {
			log.Print(err)
		}

		c.JSON(200, gin.H{
			"message": fmt.Sprintf("registration-id: %s entered the parking lot using barrier %s", reg, bID),
		})
	})

	r.POST("/barrier/:barrierID/out/:registration", func(c *gin.Context) {
		reg := c.Param("registration")
		bID := c.Param("barrierID")
		if err := bc.Validate(bID); err != nil {
			c.JSON(404, gin.H{
				"error": err,
			})
		}

		err := bc.Out(bID, reg)
		if err != nil {
			log.Print(err)
		}

		c.JSON(200, gin.H{
			"message": fmt.Sprintf("registration-id: %s exited the parking lot using barrier %s", reg, bID),
		})
	})

	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
