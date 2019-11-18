package zmq

import (
	"github.com/hyperledger/sawtooth-sdk-go/protobuf/client_batch_submit_pb2"
	"github.com/hyperledger/sawtooth-sdk-go/protobuf/validator_pb2"
	"github.com/taekion-org/sawtooth-client-sdk-go/transport/types"
)

// GetBatchStatus returns the status for a single batch.
func (self *SawtoothClientTransportZmq) GetBatchStatus(batchId string, wait int) (types.BatchStatus, error) {
	statusMap, err := self.GetBatchStatusMultiple([]string{batchId}, wait)
	if err != nil {
		return "", err
	}

	return statusMap[batchId], nil
}

// GetBatchStatusMultiple returns the statuses for a list of batches.
func (self *SawtoothClientTransportZmq) GetBatchStatusMultiple(batchIds []string, wait int) (map[string]types.BatchStatus, error) {
	// Set up the request
	t := validator_pb2.Message_CLIENT_BATCH_STATUS_REQUEST
	request := client_batch_submit_pb2.ClientBatchStatusRequest{
		BatchIds:             batchIds,
	}
	if wait > 0 {
		request.Wait = true
		request.Timeout = uint32(wait)
	}

	// Send the request and get the response
	var response client_batch_submit_pb2.ClientBatchStatusResponse
	err := self.doZmqRequest(t, &request, &response)
	if err != nil {
		return nil, err
	}

	// Create the result map from the returned data
	resultMap := make(map[string]types.BatchStatus, len(batchIds))
	for _, result := range response.BatchStatuses {
		statusString := result.Status.String()
		resultMap[result.BatchId] = types.BatchStatus(statusString)
	}

	return resultMap, nil
}
