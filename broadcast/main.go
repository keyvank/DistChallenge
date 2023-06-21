package main

import (
	"encoding/json"
	"log"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

func getNode() *maelstrom.Node {
	n := maelstrom.NewNode()
	var messages []int

	n.Handle("broadcast", func(msg maelstrom.Message) error {
		var body map[string]any
		if err := json.Unmarshal(msg.Body, &body); err != nil {
			return err
		}
		messages = append(messages, int(body["message"].(float64)))
		var resp = map[string]any{
			"type": "broadcast_ok",
		}
		return n.Reply(msg, resp)
	})

	n.Handle("read", func(msg maelstrom.Message) error {
		var body map[string]any
		if err := json.Unmarshal(msg.Body, &body); err != nil {
			return err
		}
		body["type"] = "read_ok"
		var resp = map[string]any{
			"type":     "read_ok",
			"messages": messages,
		}
		return n.Reply(msg, resp)
	})

	n.Handle("topology", func(msg maelstrom.Message) error {
		var body map[string]any
		if err := json.Unmarshal(msg.Body, &body); err != nil {
			return err
		}
		var resp = map[string]any{
			"type": "topology_ok",
		}
		return n.Reply(msg, resp)
	})

	return n
}

func main() {
	n := getNode()
	if err := n.Run(); err != nil {
		log.Fatal(err)
	}
}
