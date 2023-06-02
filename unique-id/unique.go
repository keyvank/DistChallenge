package unique

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

var nodeId string
var counter = 0

func GetNode() *maelstrom.Node {
	n := maelstrom.NewNode()
	n.Handle("init", func(msg maelstrom.Message) error {
		var body map[string]any
		if err := json.Unmarshal(msg.Body, &body); err != nil {
			return err
		}

		// Update the message type to return back.
		body["type"] = "init_ok"

		initNodeId, exists := body["node_id"]
		if !exists {
			return fmt.Errorf("no id in init message")
		}

		nodeId = initNodeId.(string)

		return n.Reply(msg, body)
	})

	n.Handle("generate", func(msg maelstrom.Message) error {
		var body map[string]any
		if err := json.Unmarshal(msg.Body, &body); err != nil {
			return err
		}

		// Update the message type to return back.
		body["type"] = "generate_ok"

		epoch := time.Now().UnixMilli()
		body["id"] = strconv.FormatInt(epoch, 10) + "-" + nodeId + "-" + strconv.FormatInt(int64(counter), 10)

		counter++

		return n.Reply(msg, body)
	})

	return n
}
