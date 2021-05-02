package api

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"

	"github.com/RicardoLinck/parking-lot-api/barrier"
)

func ConfigureEndpoints(bc *barrier.BarrierConfig) *gin.Engine {
	r := gin.Default()
	r.POST("/barrier/:barrierID/in/:registration", func(c *gin.Context) {
		reg := c.Param("registration")
		bID := c.Param("barrierID")
		if err := bc.Validate(bID); err != nil {
			c.JSON(404, gin.H{
				"error": err.Error(),
			})
			return
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
				"error": err.Error(),
			})
			return
		}

		err := bc.Out(bID, reg)
		if err != nil {
			log.Print(err)
		}

		c.JSON(200, gin.H{
			"message": fmt.Sprintf("registration-id: %s exited the parking lot using barrier %s", reg, bID),
		})
	})

	r.GET("/barrier/:barrierID/logs", func(c *gin.Context) {
		bID := c.Param("barrierID")
		if err := bc.Validate(bID); err != nil {
			c.JSON(404, gin.H{
				"error": err.Error(),
			})
			return
		}

		logs, err := bc.Logs(bID)
		if err != nil {
			c.JSON(500, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(200, logs)
	})
	return r
}
