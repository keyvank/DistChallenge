package main

import (
  "fmt"
  "encoding/json"
  "log"
  maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

func main() {
  n := maelstrom.NewNode()
  /*n.Handle("init", func(msg maelstrom.Message) error {
      // Unmarshal the message body as an loosely-typed map.
      var body map[string]any
      if err := json.Unmarshal(msg.Body, &body); err != nil {
          return err
      }

      node_id, exists := body["node_id"]
      if exists {
        n := node_id.(string)[1:2]
        rnd, err := strconv.ParseInt(n, 10, 64)
        if err != nil {
          rand.Seed(rnd)
        } else {
          return err
        }
      } else {
        return fmt.Errorf("no!!!!!!!!!!!!!!!!")
      }

      body["type"] = "init_ok"

      // Echo the original message back with the updated message type.
      return n.Reply(msg, body)
  })*/
  n.Handle("generate", func(msg maelstrom.Message) error {
      // Unmarshal the message body as an loosely-typed map.
      var body map[string]any
      if err := json.Unmarshal(msg.Body, &body); err != nil {
          return err
      }

      node_id, exists := body["node_id"]
      if !exists {
        return fmt.Errorf("no!!!!!!!!!!!!!!!!")
      }

      // Update the message type to return back.
      body["type"] = "generate_ok"
      msg_id, exists := body["msg_id"]
      if exists {
        body["id"] = int64(msg_id.(float64))
      } else {
        return fmt.Errorf("no!!!!!!!!!!!!!!!!")
      }

      // Echo the original message back with the updated message type.
      return n.Reply(msg, body)
  })
  if err := n.Run(); err != nil {
      log.Fatal(err)
  }

}
