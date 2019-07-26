package sawtooth_client_sdk_go

import (
	"encoding/hex"
	"github.com/golang/protobuf/proto"
	"github.com/hyperledger/sawtooth-sdk-go/protobuf/batch_pb2"
	"github.com/hyperledger/sawtooth-sdk-go/protobuf/transaction_pb2"
	"strings"
	"time"
)

// CreateTransaction constructs a single transaction from the provided payload.
func (self *SawtoothClient) CreateTransaction(payload interface{}) (*transaction_pb2.Transaction, error) {
	payloadEncoded, err := self.ClientImpl.EncodePayload(payload)
	if err != nil {
		return nil, err
	}

	headerPB := &transaction_pb2.TransactionHeader{
		SignerPublicKey:  self.Signer.GetPublicKey().AsHex(),
		FamilyName:       self.ClientImpl.GetFamilyName(),
		FamilyVersion:    self.ClientImpl.GetFamilyVersion(),
		Inputs:           self.ClientImpl.GetPayloadInputAddresses(payload),
		Outputs:          self.ClientImpl.GetPayloadOutputAddresses(payload),
		Dependencies:     []string{},
		PayloadSha512:    HexdigestByte(payloadEncoded),
		BatcherPublicKey: self.Signer.GetPublicKey().AsHex(),
		Nonce:            Hexdigest(time.Now().String()),
	}
	header, err := proto.Marshal(headerPB)
	if err != nil {
		return nil, err
	}

	signature := strings.ToLower(hex.EncodeToString(self.Signer.Sign(header)))

	transaction := &transaction_pb2.Transaction{
		Header:          header,
		Payload:         payloadEncoded,
		HeaderSignature: signature,
	}
	return transaction, nil
}

// CreateBatch constructs a batch from a list of payloads.
func (self *SawtoothClient) CreateBatch(transactions []*transaction_pb2.Transaction) (*batch_pb2.Batch, error) {
	var transactionSignatures []string
	for _, transaction := range transactions {
		transactionSignatures = append(transactionSignatures, transaction.HeaderSignature)
	}

	headerPB := batch_pb2.BatchHeader{
		SignerPublicKey: self.Signer.GetPublicKey().AsHex(),
		TransactionIds:  transactionSignatures,
	}
	header, err := proto.Marshal(&headerPB)
	if err != nil {
		return nil, err
	}

	signature := strings.ToLower(hex.EncodeToString(self.Signer.Sign(header)))

	batch := &batch_pb2.Batch{
		Header:          header,
		Transactions:    transactions,
		HeaderSignature: signature,
	}

	return batch, nil
}

// CreateBatchList constructs a batch list from one or more batches.
func (self *SawtoothClient) CreateBatchList(batches []*batch_pb2.Batch) (*batch_pb2.BatchList, error) {
	batchList := &batch_pb2.BatchList{
		Batches: batches,
	}

	return batchList, nil
}
