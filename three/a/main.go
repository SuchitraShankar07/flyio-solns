package main

import(
	"encoding/json"
	"log"
	"os"
	"fmt"
	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)
func main(){
	n:=maelstrom.NewNode()
	/*
	broadcast - message body:
	{"type": "broadcast",
	"message": 1000}
	*/
	var msgs []interface{}
	n.Handle("broadcast", func(message maelstrom.Message) error{

		var body map[string]interface{}
		if err:=json.Unmarshal(message.Body, &body); err != nil {
			fmt.Println("ERROR is ", err)
            return err
		}
		msgs=append(msgs, body["message"])
		returnBody:=make(map[string]interface{})
		returnBody["type"]="broadcast_ok"
		return n.Reply(message, returnBody)
	})

/*This message requests that a node return all values that it has seen.
Your node will receive a request message body that looks like this:
{
  "type": "read"
}
In response, it should return a read_ok message with a list of values it has seen:
{
  "type": "read_ok",
  "messages": [1, 8, 72, 25]
}
The order of the returned values does not matter.*/
n.Handle("read", func(msg maelstrom.Message) error{
	var body map[string]any
		if err := json.Unmarshal(msg.Body, &body); err != nil {
			return err
		}
        return n.Reply(msg, map[string]interface{}{
            "type": "read_ok",
            "messages": msgs,
        })
    })
	/*This message requests that a node return the current topology of the network.
Your node will receive a topology message body that looks like this:
{
  "type": "topology"
}
   In response, it should return a topology_ok message with a list of nodes in the network:
{
   "type": "topology_ok",
   "nodes": ["node1", "node2", "node3"]
}
   */
   
n.Handle("topology", func(message maelstrom.Message) error{
	var body map[string]interface{}
	    if err:=json.Unmarshal(message.Body, &body); err!= nil {
            fmt.Println("ERROR is ", err)
            return err
        }
		fmt.Fprintf(os.Stderr, "Get to return body\n")
		return n.Reply(message, map[string]string{
            "type": "topology_ok",
        })
})

	if err:=n.Run(); err!= nil{
		log.Fatal(err)
	}
}