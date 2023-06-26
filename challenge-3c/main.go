package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

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
		for err := n.Reply(msg, resp); err == nil; {
			continue
		}
		return nil
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

func broadcast(neighbours []string, message int, n *maelstrom.Node) error {
	for _, nodeId := range neighbours {
		body := map[string]any{
			"type":    "broadcast",
			"message": message,
		}
		sucessful := false
		for !sucessful {
			bacgroundCtx := context.Background()
			ctx, err := context.WithTimeout(bacgroundCtx, time.Second*5)
			if err != nil {
				return fmt.Errorf(">>>>> context error")
			}
			resp, timeoutErr := n.SyncRPC(ctx, nodeId, body)
			if timeoutErr != nil {
				return fmt.Errorf(">>>>> timeout error")
			}
			errCode := resp.RPCError().Code
			if errCode == 0 || errCode == 13 {
				return fmt.Errorf(">>>>> timeout error")
			} else {
				return fmt.Errorf(">>>>> timeout error: %v", errCode)
			}
		}
	}
	return nil
}

func contain(item int, list []int) bool {
	for _, value := range list {
		if value == item {
			return true
		}
	}
	return false
}
