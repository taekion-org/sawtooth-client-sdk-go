// Package transport provides a generalized interface called SawtoothClientTransport.
// A transport is an interface to communicate with Sawtooth.
// Currently, an implementation is provided that communicates with the Sawtooth REST API,
// but a ZMQ implementation that can connect directly to a validator is in the works.
package transport

import (
	"github.com/hyperledger/sawtooth-sdk-go/protobuf/batch_pb2"
	"github.com/taekion-org/sawtooth-client-sdk-go/transport/types"
)

// SawtoothClientTransport is an interface that represents a transport interface to Sawtooth.
type SawtoothClientTransport interface {
	// Methods to retrieve and submit batches.
	GetBatch(batchId string) (*types.Batch, error)
	GetBatchIterator(fetch int, reverse bool) types.BatchIterator
	GetBatchStatus(batchId string, wait int) (types.BatchStatus, error)
	GetBatchStatusMultiple(batchIds []string, wait int) (map[string]types.BatchStatus, error)
	SubmitBatchList(batchList *batch_pb2.BatchList) error

	// Methods to retrieve blocks.
	GetBlock(blockId string) (*types.Block, error)
	GetBlockIterator(fetch int, reverse bool) types.BlockIterator

	// Methods to retrieve transactions.
	GetTransaction(transactionId string) (*types.Transaction, error)
	GetTransactionIterator(fetch int, reverse bool) types.TransactionIterator

	// Methods to retrieve state.
	GetState(address string) (*types.State, error)
	GetStateAtHead(address string, head string) (*types.State, error)
	GetStateIterator(addressPrefix string, fetch int, reverse bool) types.StateIterator
}
