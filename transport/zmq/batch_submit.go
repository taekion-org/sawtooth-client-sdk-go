package zmq

import (
	"github.com/hyperledger/sawtooth-sdk-go/protobuf/batch_pb2"
	"github.com/hyperledger/sawtooth-sdk-go/protobuf/client_batch_submit_pb2"
	"github.com/hyperledger/sawtooth-sdk-go/protobuf/validator_pb2"
)

// SubmitBatchList submits a batch list to Sawtooth. The batch list must be in the form of a
// batch_pb2.BatchList protobuf and be prepared appropriately (all required fields and signatures
// populated.
func (self *SawtoothClientTransportZmq) SubmitBatchList(batchList *batch_pb2.BatchList) error {
	// Set up the request
	t := validator_pb2.Message_CLIENT_BATCH_SUBMIT_REQUEST
	request := client_batch_submit_pb2.ClientBatchSubmitRequest{
		Batches: batchList.Batches,
	}

	// Send the request and get the response
	var response client_batch_submit_pb2.ClientBatchSubmitResponse
	err := self.doZmqRequest(t, &request, &response)
	if err != nil {
		return err
	}

	return nil
}
