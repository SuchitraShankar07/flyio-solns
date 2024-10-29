package main

import (
    "encoding/json"
    "log"

    maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

func main() {
    // Create a new Maelstrom node
    n := maelstrom.NewNode()

    // Register handler for "echo" messages
    n.Handle("echo", func(msg maelstrom.Message) error {
        // Parse the message body into a generic map
        var body map[string]any
        if err := json.Unmarshal(msg.Body, &body); err != nil {
            return err
        }

        // Modify the message type to "echo_ok"
        body["type"] = "echo_ok"

        // Send the modified message back to the original sender
        return n.Reply(msg, body)
    })

    // Start the node to listen for incoming messages
    if err := n.Run(); err != nil {
        log.Fatal(err)
    }
}

