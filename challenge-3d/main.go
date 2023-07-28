package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sync"
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
	var neighboursFromTopology []string
	mutex := sync.Mutex{}

	n.Handle("broadcast", func(msg maelstrom.Message) error {
		var body map[string]any
		if err := json.Unmarshal(msg.Body, &body); err != nil {
			return err
		}
		err := func() error {
			// without the lock messages can be lost
			mutex.Lock()
			defer mutex.Unlock()
			message := int(body["message"].(float64))
			if !contain(message, messages) {
				messages = append(messages, message)
				neighbours := removeString(n.NodeIDs(), n.ID())
				if msg.Src != getLeader(n) {
					err := broadcast(message, n, neighbours)
					return err
				}
			}
			return nil
		}()
		if err != nil {
			return err
		}

		_, hasId := body["msg_id"]
		if hasId {
			reply := map[string]any{
				"type": "broadcast_ok",
			}
			return n.Reply(msg, reply)
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
			neighboursFromTopology = append(neighboursFromTopology, node.(string))
		}
		reply := map[string]any{
			"type": "topology_ok",
		}
		return n.Reply(msg, reply)
	})

	return n
}

func broadcast(message int, n *maelstrom.Node, neighbours []string) error {
	leader := getLeader(n)
	body := map[string]any{
		"type":    "broadcast",
		"message": message,
	}
	if n.ID() == leader {
		// if leader send to all
		for _, nodeId := range neighbours {
			go sendMsgWithRetry(n, body, nodeId)
		}
	} else {
		go sendMsgWithRetry(n, body, leader)
	}

	return nil
}

func sendMsgWithRetry(n *maelstrom.Node, body map[string]any, dest string) {
	success := false
	for !success {
		ctx, cancelFunc := context.WithTimeout(context.Background(), time.Second*2)
		defer cancelFunc()
		resp, err := n.SyncRPC(ctx, dest, body)
		if err != nil {
			errCode := resp.RPCError().Code
			if errCode == 0 || errCode == 13 {
				fmt.Fprintln(os.Stderr, "failed to send message ", body, "to node ", dest)
			}
			time.Sleep(time.Second)
			continue
		}
		success = true
	}
}
func getLeader(n *maelstrom.Node) string {
	leader := "n0"
	// for _, nodeId := range n.NodeIDs() {
	// 	if nodeId > leader {
	// 		leader = nodeId
	// 	}
	// }
	fmt.Fprintln(os.Stderr, "selected leader ", leader)
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
