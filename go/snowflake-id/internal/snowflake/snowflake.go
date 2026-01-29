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

	nowMs := time.Since(g.epoch).Milliseconds()
	if nowMs < 0 {
		return 0, ErrEpochInFuture
	}

	if nowMs < g.lastMs {
		return 0, ErrClockWentBack
	}

	if nowMs == g.lastMs {
		g.seq = (g.seq + 1) & maxSeq
		if g.seq == 0 {
			for nowMs <= g.lastMs {
				time.Sleep(200 * time.Microsecond)
				nowMs = time.Since(g.epoch).Milliseconds()
			}
		}
	} else {
		g.seq = 0
	}

	g.lastMs = nowMs
	ts := uint64(nowMs) & ((1 << timestampBits) - 1)

	return (ts << timeShift) | (uint64(g.node) << nodeShift) | uint64(g.seq), nil
}
