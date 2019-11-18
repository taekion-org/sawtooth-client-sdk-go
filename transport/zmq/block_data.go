package zmq

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/hyperledger/sawtooth-sdk-go/protobuf/client_block_pb2"
	"github.com/hyperledger/sawtooth-sdk-go/protobuf/client_list_control_pb2"
	"github.com/hyperledger/sawtooth-sdk-go/protobuf/validator_pb2"
	"github.com/taekion-org/sawtooth-client-sdk-go/transport/types"
)

func (self *SawtoothClientTransportZmq) GetBlock(blockId string) (*types.Block, error) {
	// Set up the request
	t := validator_pb2.Message_CLIENT_BLOCK_GET_BY_ID_REQUEST
	request := client_block_pb2.ClientBlockGetByIdRequest{
		BlockId: blockId,
	}

	// Send the request and get the response
	var response client_block_pb2.ClientBlockGetResponse
	err := self.doZmqRequest(t, &request, &response)
	if err != nil {
		return nil, err
	}

	// Convert the resulting batch into our data type
	block, err := types.BlockFromProto(response.Block)
	if err != nil {
		return nil, fmt.Errorf("Error parsing block protobuf: %s", err)
	}

	// Return the result
	return block, nil
}

type blockZmqIterator struct {
	commonZmqIterator
}

func (self *blockZmqIterator) Current() (*types.Block, error) {
	err := self.checkCurrent()
	if err != nil {
		return nil, err
	}

	data := self.getCurrent().(*types.Block)
	return data, nil
}

func (self *blockZmqIterator) BuildRequest(pagingControl *client_list_control_pb2.ClientPagingControls, sortControl []*client_list_control_pb2.ClientSortControls) (validator_pb2.Message_MessageType, proto.Message, proto.Message) {
	t := validator_pb2.Message_CLIENT_BLOCK_LIST_REQUEST
	request := client_block_pb2.ClientBlockListRequest{Paging: pagingControl, Sorting: sortControl}
	response := client_block_pb2.ClientBlockListResponse{}
	return t, &request, &response
}

func (self *blockZmqIterator) ParseProto(message proto.Message) ([]interface{}, error) {
	response := message.(*client_block_pb2.ClientBlockListResponse)
	result := make([]interface{}, len(response.Blocks))

	for i, item := range response.Blocks {
		block, err := types.BlockFromProto(item)
		if err != nil {
			return nil, err
		}
		result[i] = block
	}

	return result, nil
}

// GetBlockIterator returns a types.BlockIterator that can iterate over all blocks.
func (self *SawtoothClientTransportZmq) GetBlockIterator(fetch int, reverse bool) types.BlockIterator {
	pagingControl := &client_list_control_pb2.ClientPagingControls{
		Limit: int32(fetch),
	}

	sortControl := []*client_list_control_pb2.ClientSortControls{
		{
			Keys: []string{"block_num"},
			Reverse: reverse,
		},
	}

	iterator := &blockZmqIterator{}
	iterator.commonZmqIterator = *NewCommonZmqIterator(self, pagingControl, sortControl, iterator)

	return iterator
}
