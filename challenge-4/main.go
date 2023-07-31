package main

import (
	"context"
	"encoding/json"
	"errors"
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
	kv := maelstrom.NewSeqKV(n)
	keyName := "g"

	n.Handle("read", func(msg maelstrom.Message) error {
		ctx := context.Background()

		success := false
		var val int

		for !success {
			finalVal, err := kv.ReadInt(ctx, keyName)
			if err != nil {
				fmt.Fprintln(os.Stderr, "failed read")
				if errors.Is(err, &maelstrom.RPCError{}) {
					if err.(*maelstrom.RPCError).Code == maelstrom.KeyDoesNotExist {
						success = true
					} else {
						time.Sleep(time.Millisecond * 10)
						continue
					}
				}
			}
			err = kv.CompareAndSwap(ctx, keyName, val, val, true)
			if err == nil {
				success = true
			}
			time.Sleep(time.Millisecond * 10)
			val = finalVal
		}

		reply := map[string]any{
			"type":  "read_ok",
			"value": val,
		}
		return n.Reply(msg, reply)
	})

	n.Handle("add", func(msg maelstrom.Message) error {
		var body map[string]any
		if err := json.Unmarshal(msg.Body, &body); err != nil {
			return err
		}
		var delta int
		rawDelta, exists := body["delta"]
		if exists {
			delta = int(rawDelta.(float64))
		}

		ctx := context.Background()

		sucess := false

		for !sucess {
			prev_val, err := kv.ReadInt(ctx, keyName)
			// fmt.Fprintln(os.Stderr, "adding the value to counter", delta)
			err = kv.CompareAndSwap(ctx, keyName, prev_val, prev_val+delta, true)
			if err == nil {
				sucess = true
			}
			time.Sleep(time.Millisecond * 10)
		}

		reply := map[string]any{
			"type": "add_ok",
		}

		return n.Reply(msg, reply)

	})

	return n
}
