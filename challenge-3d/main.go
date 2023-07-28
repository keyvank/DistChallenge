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
		message := int(body["message"].(float64))
		_, exists := body["from_node"]
		if !contain(message, messages) {
			messages = append(messages, message)
		}
		if !exists {
			err := broadcast(message, n)
			if err != nil {
				return err
			}
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

func broadcast(message int, n *maelstrom.Node) error {
	allNodes := n.NodeIDs()
	neighbours := removeString(allNodes, n.ID())
	// leader := getLeader(n)
	for _, nodeId := range neighbours {
		body := map[string]any{
			"type":      "broadcast",
			"message":   message,
			"from_node": true,
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
						fmt.Fprintln(os.Stderr, "failed to send message ", message, "to node ", nodeId)
					}
					time.Sleep(time.Second)
					continue
				}
				fmt.Fprintln(os.Stderr, "successfully sent message ", message, "to node ", nodeId)
				success = true
			}
		}(nodeId)
		fmt.Fprintln(os.Stderr, time.Now().UnixNano()/int64(time.Millisecond), "sent message to neighbour from ", n.ID())
	}
	return nil
}

func getLeader(n *maelstrom.Node) string {
	leader := ""
	for _, nodeId := range n.NodeIDs() {
		if nodeId > leader {
			leader = nodeId
		}
	}
	return leader
}

func removeString(slice []string, s string) []string {
	for i, item := range slice {
		if item == s {
			return append(slice[:i], slice[i+1:]...)
		}
	}
	return slice
}

func contain(item int, list []int) bool {
	for _, value := range list {
		if value == item {
			return true
		}
	}
	return false
}
