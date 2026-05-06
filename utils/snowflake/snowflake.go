package snowflake

import (
	"sync"
	"time"
)

const (
	epochMillis  int64 = 1704067200000
	nodeBits           = 10
	sequenceBits       = 12
	maxNode            = -1 ^ (-1 << nodeBits)
	maxSequence        = -1 ^ (-1 << sequenceBits)
	nodeShift          = sequenceBits
	timeShift          = sequenceBits + nodeBits
)

type Node struct {
	mu        sync.Mutex
	nodeID    int64
	lastMilli int64
	sequence  int64
}

func NewNode(nodeID int64) *Node {
	if nodeID < 0 || nodeID > maxNode {
		nodeID = 1
	}
	return &Node{nodeID: nodeID}
}

func (n *Node) NextID() int64 {
	n.mu.Lock()
	defer n.mu.Unlock()

	now := time.Now().UnixMilli()
	if now == n.lastMilli {
		n.sequence = (n.sequence + 1) & maxSequence
		if n.sequence == 0 {
			for now <= n.lastMilli {
				now = time.Now().UnixMilli()
			}
		}
	} else {
		n.sequence = 0
	}
	n.lastMilli = now
	return ((now - epochMillis) << timeShift) | (n.nodeID << nodeShift) | n.sequence
}
