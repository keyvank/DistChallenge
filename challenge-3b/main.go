package main

import (
	"encoding/json"
	"errors"
	"log"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

func main() {
	n := getNode()
	if err := n.Run(); err != nil {
		log.Fatal(err)
	}
}

func getNode() *maelstrom.Node {
	n := maelstrom.NewNode()
	var messages []int
	var neighbours []string

	n.Handle("broadcast", func(msg maelstrom.Message) error {
		var body map[string]any
		if err := json.Unmarshal(msg.Body, &body); err != nil {
			return err
		}
		message := int(body["message"].(float64))
		if !contain(message, messages) {
			messages = append(messages, message)
			broadcast(neighbours, message, n)
		}
		resp := map[string]any{
			"type": "broadcast_ok",
		}
		return n.Reply(msg, resp)
	})

	n.Handle("read", func(msg maelstrom.Message) error {
		var body map[string]any
		if err := json.Unmarshal(msg.Body, &body); err != nil {
			return err
		}
		reply := map[string]any{
			"type":     "read_ok",
			"messages": messages,
		}
		return n.Reply(msg, reply)
	})

	n.Handle("topology", func(msg maelstrom.Message) error {
		var body map[string]any
		if err := json.Unmarshal(msg.Body, &body); err != nil {
			return err
		}
		for _, node := range body["topology"].(map[string]any)[n.ID()].([]any) {
			neighbours = append(neighbours, node.(string))
		}
		reply := map[string]any{
			"type": "topology_ok",
		}
		return n.Reply(msg, reply)
	})

	return n
}

func broadcast(neighbours []string, message int, n *maelstrom.Node) {
	for _, nodeId := range neighbours {
		reply := map[string]any{
			"type":    "broadcast",
			"message": message,
		}
		n.RPC(nodeId, reply, func(msg maelstrom.Message) error {
			var body map[string]any
			if err := json.Unmarshal(msg.Body, &body); err != nil {
				return err
			}
			if body["type"].(string) != "broadcast_ok" {
				return errors.New("broadcast failed")
			}
			return nil
		})
	}
}

func contain(item int, list []int) bool {
	for _, value := range list {
		if value == item {
			return true
		}
	}
	return false
}
