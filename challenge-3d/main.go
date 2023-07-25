package main

import (
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
		// if !contain(message, messages) {
		messages = append(messages, message)
		err := broadcast(message, n)
		if err != nil {
			return err
		}
		// }
		reply := map[string]any{
			"type": "broadcast_ok",
		}
		return n.Reply(msg, reply)
	})

	n.Handle("propagate", func(msg maelstrom.Message) error {
		var body map[string]any
		if err := json.Unmarshal(msg.Body, &body); err != nil {
			return err
		}

		message := int(body["message"].(float64))
		messages = append(messages, message)

		// reply := map[string]any{
		// 	"type": "propagate_ok",
		// }
		// return n.Reply(msg, reply)
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

func broadcast(message int, n *maelstrom.Node) error {
	neighbours := n.NodeIDs()
	for _, nodeId := range neighbours {
		if n.ID() == nodeId {
			continue
		}
		body := map[string]any{
			"type":    "propagate",
			"message": message,
		}
		n.Send(nodeId, body)
		fmt.Fprintln(os.Stderr, time.Now().UnixNano()/int64(time.Millisecond), "sent message to neighbour from ", n.ID())
		// go func(nodeId string) {
		// 	success := false
		// 	for !success {
		// 		ctx, cancelFunc := context.WithTimeout(context.Background(), time.Second*2)
		// 		defer cancelFunc()
		// 		resp, err := n.SyncRPC(ctx, nodeId, body)
		// 		if err != nil {
		// 			errCode := resp.RPCError().Code
		// 			if errCode == 0 || errCode == 13 {
		// 				fmt.Fprintln(os.Stderr, "failed to send message ", message, "to node ", nodeId)
		// 			}
		// 			time.Sleep(time.Second)
		// 			continue
		// 		}
		// 		fmt.Fprintln(os.Stderr, "successfully sent message ", message, "to node ", nodeId)
		// 		success = true
		// 	}
		// }(nodeId)
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
