package main

import (
	"github.com/RicardoLinck/parking-lot-api/api"
	"github.com/RicardoLinck/parking-lot-api/barrier"
)

func main() {
	bc := barrier.NewBarrierConfig("./logs/barriers")
	api.ConfigureEndpoints(bc).Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
