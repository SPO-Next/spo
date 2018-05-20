package webrpc

import (
	"sort"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/spo-next/spo/src/cipher"
	"github.com/spo-next/spo/src/coin"
	"github.com/spo-next/spo/src/testutil"
	"github.com/spo-next/spo/src/util/logging"
	"github.com/spo-next/spo/src/visor"
)

var (
	log = logging.MustGetLogger("webrpc_test")
)

// Tests are setup as subtests, to retain a single *WebRPC instance for scaffolding
// https://blog.golang.org/subtests
func TestClient(t *testing.T) {
	s := setupWebRPC(t)
	errC := make(chan error, 1)

	go func() {
		errC <- s.Run()
	}()

	time.Sleep(time.Millisecond * 100) // give s.Run() enough time to start

	defer func() {
		err := s.Shutdown()
		require.NoError(t, err)
		require.NoError(t, <-errC)
	}()

	c := &Client{
		Addr: s.Addr,
	}

	testFuncs := []struct {
		n string
		f func(t *testing.T, c *Client, s *WebRPC, gw *fakeGateway)
	}{
		{"get unspent outputs", testClientGetUnspentOutputs},
		{"inject transaction", testClientInjectTransaction},
		{"get status", testClientGetStatus},
		{"get transaction by id", testClientGetTransactionByID},
		{"get address uxouts", testClientGetAddressUxOuts},
		{"get blocks", testClientGetBlocks},
		{"get blocks by seq", testClientGetBlocksBySeq},
		{"get last block", testClientGetLastBlocks},
	}

	for _, f := range testFuncs {
		t.Run(f.n, func(t *testing.T) {
			f.f(t, c, s, s.Gateway.(*fakeGateway))
		})
	}
}

func testClientGetUnspentOutputs(t *testing.T, c *Client, s *WebRPC, gw *fakeGateway) {
	headTime := uint64(time.Now().UTC().Unix())
	uxouts := make([]coin.UxOut, 5)
	addrs := make([]cipher.Address, 5)
	rbOutputs := make(visor.ReadableOutputs, 5)
	for i := 0; i < 5; i++ {
		addrs[i] = testutil.MakeAddress()
		uxouts[i] = coin.UxOut{}
		uxouts[i].Body.Address = addrs[i]
		rbOut, err := visor.NewReadableOutput(headTime, uxouts[i])
		require.NoError(t, err)
		rbOutputs[i] = rbOut
	}

	s.Gateway = &fakeGateway{
		uxouts: uxouts,
	}

	defer func() {
		s.Gateway = gw
	}()

	reqAddrs := []string{addrs[0].String(), addrs[1].String()}

	outputs, err := c.GetUnspentOutputs(reqAddrs)
	require.NoError(t, err)
	require.Len(t, outputs.Outputs.HeadOutputs, 2)
	require.Len(t, outputs.Outputs.IncomingOutputs, 0)
	require.Len(t, outputs.Outputs.OutgoingOutputs, 0)

	// GetUnspentOutputs sorts outputs by most recent time first, then by hash
	expectedOutputs := rbOutputs[:2]
	sort.Slice(expectedOutputs, func(i, j int) bool {
		if expectedOutputs[i].Time == expectedOutputs[j].Time {
			return strings.Compare(expectedOutputs[i].Hash, expectedOutputs[j].Hash) < 1
		}

		return expectedOutputs[i].Time > expectedOutputs[j].Time
	})

	require.Equal(t, rbOutputs[:2], outputs.Outputs.HeadOutputs)

	// Invalid address
	_, err = c.GetUnspentOutputs([]string{"invalid-address-foo"})
	require.Error(t, err)
	require.Equal(t, "invalid address: invalid-address-foo [code: -32602]", err.Error())
}

func testClientInjectTransaction(t *testing.T, c *Client, s *WebRPC, gw *fakeGateway) {
	gw.injectRawTxMap = map[string]bool{
		rawTxID: true,
	}
	require.Empty(t, gw.injectedTransactions)

	txID, err := c.InjectTransactionString(rawTxStr)
	require.NoError(t, err)
	require.NotEmpty(t, txID)

	log.Println(gw.injectedTransactions)
	require.Len(t, gw.injectedTransactions, 1)
	require.Contains(t, gw.injectedTransactions, rawTxID)
}

func testClientGetStatus(t *testing.T, c *Client, s *WebRPC, gw *fakeGateway) {
	status, err := c.GetStatus()
	require.NoError(t, err)
	// values derived from hardcoded `blockString`
	require.Equal(t, &StatusResult{
		Running:            true,
		BlockNum:           455,
		LastBlockHash:      "",
		TimeSinceLastBlock: "18446744072232256374s",
	}, status)
}

func testClientGetTransactionByID(t *testing.T, c *Client, s *WebRPC, gw *fakeGateway) {
	// Invalid txn id (not SHA256)
	txid := "foo"
	txn, err := c.GetTransactionByID(txid)
	require.Nil(t, txn)
	require.Error(t, err)

	// Valid txn id, txn does not exist
	// TODO
	txn, err = c.GetTransactionByID(rawTxID)
	require.Nil(t, txn)
	require.Error(t, err)

	// Txn exists
	gw.transactions = map[string]string{
		rawTxID: rawTxStr,
	}
	txn, err = c.GetTransactionByID(rawTxID)
	require.NoError(t, err)
	expectedTxn := decodeRawTransaction(rawTxStr)
	rbTx, err := visor.NewReadableTransaction(expectedTxn)
	require.NoError(t, err)
	require.Equal(t, &visor.TransactionResult{
		Status:      expectedTxn.Status,
		Time:        0,
		Transaction: *rbTx,
	}, txn.Transaction)
}

func testClientGetAddressUxOuts(t *testing.T, c *Client, s *WebRPC, gw *fakeGateway) {
	// UxOut tests use the autogenerated gateway mock instead of *fakeGateway
	// Temporarily swap it onto the *WebRPC server
	gatewayerMock, mockData := newUxOutMock()
	s.Gateway = gatewayerMock
	defer func() {
		s.Gateway = gw
	}()

	addr := "2kmKohJrwURrdcVtDNaWK6hLCNsWWbJhTqT"
	addrs := []string{addr}

	uxouts, err := c.GetAddressUxOuts(addrs)
	require.NoError(t, err)

	require.Equal(t, []AddrUxoutResult{{
		Address: addr,
		UxOuts:  mockData(addr),
	}}, uxouts)

	// Invalid addr
	uxouts, err = c.GetAddressUxOuts([]string{"foo"})
	require.Error(t, err)
	require.Nil(t, uxouts)
}

func testClientGetBlocks(t *testing.T, c *Client, s *WebRPC, gw *fakeGateway) {
	// blockString borrowed from block_test.go
	blocks, err := c.GetBlocks(0, 1)
	require.NoError(t, err)
	require.Equal(t, decodeBlock(blockString), blocks)
}

func testClientGetBlocksBySeq(t *testing.T, c *Client, s *WebRPC, gw *fakeGateway) {
	gatewayerMock := NewGatewayerMock()
	s.Gateway = gatewayerMock
	defer func() {
		s.Gateway = gw
	}()

	gatewayerMock.On("GetBlocksInDepth", []uint64{454}).Return(decodeBlock(blockString), nil)

	// blockString and seq borrowed from block_test.go
	var seq uint64 = 454
	blocks, err := c.GetBlocksBySeq([]uint64{seq})
	require.NotNil(t, blocks.Blocks)
	require.NoError(t, err)
	require.Equal(t, decodeBlock(blockString), blocks)
}

func testClientGetLastBlocks(t *testing.T, c *Client, s *WebRPC, gw *fakeGateway) {
	var n uint64 = 1
	blocks, err := c.GetLastBlocks(n)
	require.NoError(t, err)
	require.Len(t, blocks.Blocks, 1)
	require.Equal(t, decodeBlock(blockString), blocks)
}
