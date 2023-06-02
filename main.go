package main

import (
	"log"

	"github.com/keyvank/DistChallenge/unique-id"
)

var nodeId string
var counter = 0

func main() {
	n := unique.GetNode()
	if err := n.Run(); err != nil {
		log.Fatal(err)
	}

}
