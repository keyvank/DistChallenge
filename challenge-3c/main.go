package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
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
		fmt.Fprintln(os.Stderr, "broadcast called")
		message := int(body["message"].(float64))
		if !contain(message, messages) {
			messages = append(messages, message)
			result := broadcast(neighbours, message, n)
			fmt.Fprintln(os.Stderr, "broadcast result, ", result)
		}
		reply := map[string]any{
			"type": "broadcast_ok",
		}
		return n.Reply(msg, reply)
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
		go func(nodeId string) {
			success := false
			for !success {
				ctx, cancelFunc := context.WithTimeout(context.Background(), time.Second*2)
				defer cancelFunc()
				resp, err := n.SyncRPC(ctx, nodeId, body)
				if err != nil {
					errCode := resp.RPCError().Code
					if errCode == 0 || errCode == 13 {
						fmt.Fprintln(os.Stderr, "f to send message ", message, "to node ", nodeId)
					}
					time.Sleep(time.Second)
					continue
				}
				fmt.Fprintln(os.Stderr, "s to send message ", message, "to node ", nodeId)
				success = true
			}
		}(nodeId)
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
