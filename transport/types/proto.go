package types

import (
	"github.com/golang/protobuf/proto"
	"github.com/hyperledger/sawtooth-sdk-go/protobuf/batch_pb2"
	"github.com/hyperledger/sawtooth-sdk-go/protobuf/block_pb2"
	"github.com/hyperledger/sawtooth-sdk-go/protobuf/transaction_pb2"
)

// TransactionFromProto converts a Transaction protobuf into our own Transaction object.
func TransactionFromProto(transactionProto *transaction_pb2.Transaction) (*Transaction, error) {
	// Parse out the transaction header
	var headerProto transaction_pb2.TransactionHeader
	err := proto.Unmarshal(transactionProto.Header, &headerProto)
	if err != nil {
		return nil, err
	}

	// Construct the transaction object
	transaction := Transaction{
		Header: TransactionHeader{
			BatcherPublicKey: headerProto.BatcherPublicKey,
			Dependencies:     headerProto.Dependencies,
			FamilyName:       headerProto.FamilyName,
			FamilyVersion:    headerProto.FamilyVersion,
			Inputs:           headerProto.Inputs,
			Nonce:            headerProto.Nonce,
			Outputs:          headerProto.Outputs,
			PayloadSHA256:    headerProto.PayloadSha512,
			SignerPublicKey:  headerProto.SignerPublicKey,
		},
		HeaderSignature: transactionProto.HeaderSignature,
		Payload: transactionProto.Payload,
	}

	return &transaction, nil
}


// BatchFromProto converts a Batch protobuf into our own Batch object.
func BatchFromProto(batchProto *batch_pb2.Batch) (*Batch, error) {
	// Parse out the batch header
	var headerProto batch_pb2.BatchHeader
	err := proto.Unmarshal(batchProto.Header, &headerProto)
	if err != nil {
		return nil, err
	}

	// Construct the batch object
	batch := Batch{
		Header: BatchHeader{
			SignerPublicKey: headerProto.SignerPublicKey,
			TransactionIds:  headerProto.TransactionIds,
		},
		HeaderSignature: batchProto.HeaderSignature,
		Transactions: make([]Transaction, len(batchProto.Transactions)),
		Trace: batchProto.Trace,
	}

	// Parse out the transactions
	for i, transactionProto := range(batchProto.Transactions) {
		transaction, err := TransactionFromProto(transactionProto)
		if err != nil {
			return nil, err
		}
		batch.Transactions[i] = *transaction
	}

	return &batch, nil
}

// BlockFromProto converts a Block protobuf into our own Block object.
func BlockFromProto(blockProto *block_pb2.Block) (*Block, error) {
	// Parse out the block header
	var headerProto block_pb2.BlockHeader
	err := proto.Unmarshal(blockProto.Header, &headerProto)
	if err != nil {
		return nil, err
	}

	// Construct the block object
	block := Block{
		Header: BlockHeader{
			BatchIds:        headerProto.BatchIds,
			BlockNum:        string(headerProto.BlockNum),
			Consensus:       headerProto.Consensus,
			PreviousBlockId: headerProto.PreviousBlockId,
			SignerPublicKey: headerProto.SignerPublicKey,
			StateRootHash:   headerProto.StateRootHash,
		},
		HeaderSignature: blockProto.HeaderSignature,
		Batches: make([]Batch, len(blockProto.Batches)),
	}

	// Parse out the batches
	for i, batchProto := range(blockProto.Batches) {
		batch, err := BatchFromProto(batchProto)
		if err != nil {
			return nil, err
		}
		block.Batches[i] = *batch
	}

	return &block, nil
}
