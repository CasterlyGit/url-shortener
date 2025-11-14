package snowflake

import (
    "errors"
    "sync"
    "time"
)

const (
    epoch         = int64(1609459200000) // Custom epoch (January 1, 2021)
    nodeBits      = 10                   // Number of bits for node ID
    sequenceBits  = 12                   // Number of bits for sequence number
    nodeMax       = -1 ^ (-1 << nodeBits)
    sequenceMask  = -1 ^ (-1 << sequenceBits)
    nodeShift     = sequenceBits
    timestampShift = sequenceBits + nodeBits
)

type Node struct {
    mu        sync.Mutex
    timestamp int64
    nodeID    int64
    sequence  int64
}

func NewNode(nodeID int64) (*Node, error) {
    if nodeID < 0 || nodeID > nodeMax {
        return nil, errors.New("node ID out of range")
    }
    return &Node{
        timestamp: 0,
        nodeID:    nodeID,
        sequence:  0,
    }, nil
}

func (n *Node) Generate() int64 {
    n.mu.Lock()
    defer n.mu.Unlock()

    now := time.Now().UnixMilli()

    if now == n.timestamp {
        n.sequence = (n.sequence + 1) & sequenceMask
        if n.sequence == 0 {
            // Sequence exhausted, wait for next millisecond
            for now <= n.timestamp {
                now = time.Now().UnixMilli()
            }
        }
    } else {
        n.sequence = 0
    }

    n.timestamp = now

    return (now-epoch)<<timestampShift |
        (n.nodeID << nodeShift) |
        n.sequence
}