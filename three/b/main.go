package main

import (
	"encoding/json"
	"log"
	"sync"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)
//sync gives us the required concurrency primitives like mutex to allow us to have concurrent access to the data
type server struct {
	n      *maelstrom.Node
	nodeID string
	idsMu  sync.RWMutex
	ids    map[int]struct{}
}
/*RWMutex: a read write mutex, itself is a struct of structure
    type RWMutex struct {
	w           Mutex        // held if there are pending writers
	writerSem   uint32       // semaphore for writers to wait for completing readers
	readerSem   uint32       // semaphore for readers to wait for completing writers
	readerCount atomic.Int32 // number of pending readers
	readerWait  atomic.Int32 // number of departing readers
}

    you cannot copy RWMutex after its first use, and only one thread can write to ids, while all can read simultaneously

*/
func main() {
	n := maelstrom.NewNode()
	s := &server{n: n, nodeID: n.ID(), ids: make(map[int]struct{})}

	// Broadcast handler
	n.Handle("broadcast", func(msg maelstrom.Message) error {
		var body map[string]any
		if err := json.Unmarshal(msg.Body, &body); err != nil {
			return err
		}

		id := int(body["message"].(float64))
		s.idsMu.Lock()
		if _, exists := s.ids[id]; exists {
			s.idsMu.Unlock()
			return nil
		}
		s.ids[id] = struct{}{}
		s.idsMu.Unlock()

		if err := s.broadcast(msg.Src, body); err != nil {
			return err
		}

		return s.n.Reply(msg, map[string]any{
			"type": "broadcast_ok",
		})
	})

	n.Handle("read", func(msg maelstrom.Message) error {
		ids := s.getAllIDs()

		return s.n.Reply(msg, map[string]any{
			"type":     "read_ok",
			"messages": ids,
		})
	})

	n.Handle("topology", func(msg maelstrom.Message) error {
		return n.Reply(msg, map[string]string{
			"type": "topology_ok",
		})
	})

	if err := n.Run(); err != nil {
		log.Fatal(err)
	}
}

func (s *server) broadcast(src string, body map[string]any) error {
	for _, dst := range s.n.NodeIDs() {
		if dst == src || dst == s.nodeID {
			continue
		}

		go func(dst string) {
			if err := s.n.Send(dst, body); err != nil {
				panic(err)
			}
		}(dst)
	}
	return nil
}

func (s *server) getAllIDs() []int {
	s.idsMu.RLock()
	defer s.idsMu.RUnlock()

	ids := make([]int, 0, len(s.ids))
	for id := range s.ids {
		ids = append(ids, id)
	}
	return ids
}
