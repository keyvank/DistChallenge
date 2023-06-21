package main

import (
	"encoding/json"
	"log"
	"strconv"
	"time"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

func getNode() *maelstrom.Node {
	n := maelstrom.NewNode()
	var counter = 0

	n.Handle("generate", func(msg maelstrom.Message) error {
		var body map[string]any
		if err := json.Unmarshal(msg.Body, &body); err != nil {
			return err
		}
		body["type"] = "generate_ok"
		epoch := time.Now().UnixMilli()
		body["id"] = strconv.FormatInt(epoch, 10) + "-" + n.ID() + "-" + strconv.FormatInt(int64(counter), 10)
		counter++
		return n.Reply(msg, body)
	})

	return n
}

func main() {
	n := getNode()
	if err := n.Run(); err != nil {
		log.Fatal(err)
	}
}
