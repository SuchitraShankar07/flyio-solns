package main

import (
	"encoding/json"
	"log"
	"strconv"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)
//my approach to this problem involves a simple upcounter which generates unique id based on the previous node id, just increments by 1 for every new id.



func main() {
	node := maelstrom.NewNode()
	var lastId uint64

	node.Handle("generate", generator(node, &lastId))

	if err := node.Run(); err != nil {
		log.Fatalf("Node run error: %v", err)
	}
}

func generator(node *maelstrom.Node, lastId *uint64) func(msg maelstrom.Message) error {

	return func(msg maelstrom.Message) error {
		var requestBody map[string]any
		if err := json.Unmarshal(msg.Body, &requestBody); err != nil {
			return err
		}

		uniqueID := uniqueId(node.ID(), lastId)
		responseBody := map[string]any{
			"type": "generate_ok",
			"id":   uniqueID,
		}

		return node.Reply(msg, responseBody)
	}
}

func uniqueId(nodeID string, lastId *uint64) string {
	*lastId++
	return nodeID + "_" + strconv.FormatUint(*lastId, 10)
}

//ideally I'd use a uuid, but since this is a high throughput system, possible race conditions/ latency/network partitioning issues could have come up. 
//perhaps combining uuid generation with timestamp generation, in a hybrid approach would be better
//however simple an upcounter is, it is foolproof. also it works. so dont mess with it.
// it is extremely rare to have the same uuid consecutively (For example, the number of random version-4 UUIDs which need to be generated in order to have a 50% probability of at least one collision is 2.71 quintillion) 
// read (https://en.wikipedia.org/wiki/Universally_unique_identifier)