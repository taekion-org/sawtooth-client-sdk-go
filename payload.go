package sawtooth_client_sdk_go

import (
	"fmt"
	"github.com/hyperledger/sawtooth-sdk-go/protobuf/batch_pb2"
	"github.com/hyperledger/sawtooth-sdk-go/protobuf/transaction_pb2"
	"github.com/taekion-org/sawtooth-client-sdk-go/transport/types"
	"time"
)

// ExecutePayload submits a single transaction to the blockchain and returns the batch id.
func (self *SawtoothClient) ExecutePayload(payload interface{}) (string, error) {
	payloads := []interface{}{payload}
	return self.ExecutePayloadBatch(payloads)
}

// ExecutePayloadSync submits a single transaction to the blockchain and waits for commit.
func (self *SawtoothClient) ExecutePayloadSync(payload interface{}, timeout int, pollInterval int) error {
	payloads := []interface{}{payload}
	return self.ExecutePayloadBatchSync(payloads, timeout, pollInterval)
}

// ExecutePayload submits a list of transactions to the blockchain (as a single batch) and returns the batch id.
func (self *SawtoothClient) ExecutePayloadBatch(payloads []interface{}) (string, error) {
	transactions := make([]*transaction_pb2.Transaction, len(payloads))

	for i, payload := range payloads {
		transaction, err := self.CreateTransaction(payload)
		if err != nil {
			return "", fmt.Errorf("Error while creating transaction for payload %d: %s", i, err)
		}
		transactions[i] = transaction
	}

	batch, err := self.CreateBatch(transactions)
	if err != nil {
		return "", err
	}
	batchId := batch.HeaderSignature

	batches := []*batch_pb2.Batch{batch}
	batchList, err := self.CreateBatchList(batches)
	if err != nil {
		return "", err
	}

	err = self.Transport.SubmitBatchList(batchList)
	if err != nil {
		return "", err
	}

	return batchId, nil
}

// ExecutePayload submits a list of transactions to the blockchain (as a single batch) and waits for commit.
func (self *SawtoothClient) ExecutePayloadBatchSync(payloads []interface{}, timeout int, pollInterval int) error {
	// Execute the payload
	batchId, err := self.ExecutePayloadBatch(payloads)
	if err != nil {
		return err
	}

	// Poll for the payload to be executed (batch committed)
	success, err := self.WaitBatch(batchId, timeout, pollInterval)
	if err != nil {
		return err
	}

	if !success {
		return fmt.Errorf("Batch not commited as of configured timeout")
	}

	return nil
}

// WaitBatch performs a polling wait for a particular batch.
func (self *SawtoothClient) WaitBatch(batchId string, timeout int, pollInterval int) (bool, error) {
	startTime := time.Now().Unix()
	waitTime := 0

	for {
		status, err := self.Transport.GetBatchStatus(batchId, pollInterval)
		if err != nil {
			return false, err
		}

		waitTime = int(time.Now().Unix() - startTime)

		switch status {
		case types.BATCH_STATUS_PENDING, types.BATCH_STATUS_UNKNOWN:
			if (timeout == 0) || (timeout != 0 && waitTime < timeout) {
				continue
			} else {
				return false, nil
			}
		case types.BATCH_STATUS_COMMITTED:
			return true, nil
		case types.BATCH_STATUS_INVALID:
			return false, fmt.Errorf("Batch %s is in status %s", batchId, status)
		}
	}

	return false, nil
}
