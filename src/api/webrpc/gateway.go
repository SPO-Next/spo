package webrpc

import (
	"github.com/spolabs/spo/src/cipher"
	"github.com/spolabs/spo/src/coin"
	"github.com/spolabs/spo/src/daemon"
	"github.com/spolabs/spo/src/visor"
	"github.com/spolabs/spo/src/visor/historydb"
)

//go:generate goautomock -template=testify Gatewayer

// Gatewayer provides interfaces for getting spo related info.
type Gatewayer interface {
	GetLastBlocks(num uint64) (*visor.ReadableBlocks, error)
	GetBlocks(start, end uint64) (*visor.ReadableBlocks, error)
	GetBlocksInDepth(vs []uint64) (*visor.ReadableBlocks, error)
	GetUnspentOutputs(filters ...daemon.OutputsFilter) (visor.ReadableOutputSet, error)
	GetTransaction(txid cipher.SHA256) (*visor.Transaction, error)
	InjectTransaction(tx coin.Transaction) error
	GetAddrUxOuts(addr cipher.Address) ([]*historydb.UxOutJSON, error)
	GetTimeNow() uint64
}
