package snowflake

import (
	"errors"
	"sync"
	"time"
)

const (
	// Layout: 41 bits timestamp (ms) | 10 bits node | 12 bits sequence
	timestampBits = 41
	nodeBits      = 10
	seqBits       = 12

	maxNode = (1 << nodeBits) - 1
	maxSeq  = (1 << seqBits) - 1

	nodeShift = seqBits
	timeShift = nodeBits + seqBits
)

var (
	ErrInvalidNodeID   = errors.New("invalid node id")
	ErrClockWentBack   = errors.New("clock moved backwards")
	ErrEpochInFuture   = errors.New("epoch is in the future")
	defaultCustomEpoch = time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
)

type Generator struct {
	mu    sync.Mutex
	epoch time.Time
	node  uint16

	lastMs int64
	seq    uint16
}

func New(nodeID uint16) (*Generator, error) {
	return NewWithEpoch(nodeID, defaultCustomEpoch)
}

func NewWithEpoch(nodeID uint16, epoch time.Time) (*Generator, error) {
	if nodeID > maxNode {
		return nil, ErrInvalidNodeID
	}
	if epoch.After(time.Now().UTC()) {
		return nil, ErrEpochInFuture
	}
	return &Generator{
		epoch:  epoch.UTC(),
		node:   nodeID,
		lastMs: -1,
	}, nil
}

func (g *Generator) NextID() (uint64, error) {
	g.mu.Lock()
	defer g.mu.Unlock()

	// Measure time in milliseconds from the custom epoch.
	nowMs := time.Since(g.epoch).Milliseconds()
	if nowMs < 0 {
		return 0, ErrEpochInFuture
	}

	// Reject IDs if the observed clock goes backwards.
	if nowMs < g.lastMs {
		return 0, ErrClockWentBack
	}

	if nowMs == g.lastMs {
		// Same millisecond: advance the per-millisecond sequence.
		g.seq = (g.seq + 1) & maxSeq
		if g.seq == 0 {
			// Sequence wrapped, so wait until time moves to the next millisecond.
			for nowMs <= g.lastMs {
				time.Sleep(200 * time.Microsecond)
				nowMs = time.Since(g.epoch).Milliseconds()
			}
		}
	} else {
		// New millisecond: start the sequence from zero again.
		g.seq = 0
	}

	// Remember the latest millisecond we used for the next call.
	g.lastMs = nowMs

	// Keep only the timestamp bits that fit in the Snowflake layout.
	ts := uint64(nowMs) & ((1 << timestampBits) - 1)

	// Pack timestamp, node id, and sequence into one 64-bit ID.
	// Example with small made-up values:
	//   timestamp = 5  -> 101
	//   node      = 3  -> 11
	//   seq       = 2  -> 10
	// After shifting into position:
	//   timestamp << timeShift  puts timestamp in the high bits
	//   uint64(g.node) << nodeShift moves node left by 12 bits, so:
	//   3        = ...000000000011
	//   3 << 12  = ...0011000000000000
	// This leaves the lowest 12 bits free for the sequence value.
	// The final ID is built by OR-ing the three parts together.
	return (ts << timeShift) | (uint64(g.node) << nodeShift) | uint64(g.seq), nil
}
