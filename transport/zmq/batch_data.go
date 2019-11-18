package zmq

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/hyperledger/sawtooth-sdk-go/protobuf/client_batch_pb2"
	"github.com/hyperledger/sawtooth-sdk-go/protobuf/client_list_control_pb2"
	"github.com/hyperledger/sawtooth-sdk-go/protobuf/validator_pb2"
	"github.com/taekion-org/sawtooth-client-sdk-go/transport/types"
)

// GetBatch returns the batch represented by batchId.
func (self *SawtoothClientTransportZmq) GetBatch(batchId string) (*types.Batch, error) {
	// Set up the request
	t := validator_pb2.Message_CLIENT_BATCH_GET_REQUEST
	request := client_batch_pb2.ClientBatchGetRequest{
		BatchId: batchId,
	}

	// Send the request and get the response
	var response client_batch_pb2.ClientBatchGetResponse
	err := self.doZmqRequest(t, &request, &response)
	if err != nil {
		return nil, err
	}

	// Convert the resulting batch into our data type
	batch, err := types.BatchFromProto(response.Batch)
	if err != nil {
		return nil, fmt.Errorf("Error parsing batch protobuf: %s", err)
	}

	// Return the result
	return batch, nil
}

type batchZmqIterator struct {
	commonZmqIterator
}

func (self *batchZmqIterator) Current() (*types.Batch, error) {
	err := self.checkCurrent()
	if err != nil {
		return nil, err
	}

	data := self.getCurrent().(*types.Batch)
	return data, nil
}

func (self *batchZmqIterator) BuildRequest(pagingControl *client_list_control_pb2.ClientPagingControls, sortControl []*client_list_control_pb2.ClientSortControls) (validator_pb2.Message_MessageType, proto.Message, proto.Message) {
	t := validator_pb2.Message_CLIENT_BATCH_LIST_REQUEST
	request := client_batch_pb2.ClientBatchListRequest{Paging: pagingControl, Sorting: sortControl}
	response := client_batch_pb2.ClientBatchListResponse{}
	return t, &request, &response
}

func (self *batchZmqIterator) ParseProto(message proto.Message) ([]interface{}, error) {
	response := message.(*client_batch_pb2.ClientBatchListResponse)
	result := make([]interface{}, len(response.Batches))

	for i, item := range response.Batches {
		batch, err := types.BatchFromProto(item)
		if err != nil {
			return nil, err
		}
		result[i] = batch
	}

	return result, nil
}

// GetBatchIterator returns a types.BatchIterator that can iterate over all batches.
func (self *SawtoothClientTransportZmq) GetBatchIterator(fetch int, reverse bool) types.BatchIterator {
	pagingControl := &client_list_control_pb2.ClientPagingControls{
		Limit: int32(fetch),
	}

	sortControl := []*client_list_control_pb2.ClientSortControls{
		{
			Keys: []string{"default"},
			Reverse: reverse,
		},
	}

	iterator := &batchZmqIterator{}
	iterator.commonZmqIterator = *NewCommonZmqIterator(self, pagingControl, sortControl, iterator)

	return iterator
}
