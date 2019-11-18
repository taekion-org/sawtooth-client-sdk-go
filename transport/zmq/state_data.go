package zmq

import (
	"github.com/golang/protobuf/proto"
	"github.com/hyperledger/sawtooth-sdk-go/protobuf/block_pb2"
	"github.com/hyperledger/sawtooth-sdk-go/protobuf/client_block_pb2"
	"github.com/hyperledger/sawtooth-sdk-go/protobuf/client_list_control_pb2"
	"github.com/hyperledger/sawtooth-sdk-go/protobuf/client_state_pb2"
	"github.com/hyperledger/sawtooth-sdk-go/protobuf/validator_pb2"
	"github.com/taekion-org/sawtooth-client-sdk-go/transport/types"
)

// headToStateRoot returns state root for the block specified.
func (self *SawtoothClientTransportZmq) headToStateRoot(blockId string) (string, error) {

	// Set up the request
	t := validator_pb2.Message_CLIENT_BLOCK_GET_BY_ID_REQUEST
	request := client_block_pb2.ClientBlockGetByIdRequest{
		BlockId: blockId,
	}

	// Send the request and get the response
	var response client_block_pb2.ClientBlockGetResponse
	err := self.doZmqRequest(t, &request, &response)
	if err != nil {
		return "", err
	}

	header, err := self.decodeBlockHeader(response.Block)
	if err != nil {
		return "", err
	}

	return header.StateRootHash, nil
}

// currentStateRoot retrieves the most recent block and returns the head (block id) and state root
func (self *SawtoothClientTransportZmq) currentStateRoot() (string, string, error) {
	// We need to get the most recent block and retrieve the state root

	// Set up the request
	t := validator_pb2.Message_CLIENT_BLOCK_LIST_REQUEST
	request := client_block_pb2.ClientBlockListRequest{
		Paging: &client_list_control_pb2.ClientPagingControls{
			Limit: 1,
		},
	}

	// Send the request and get the response
	var response client_block_pb2.ClientBlockListResponse
	err := self.doZmqRequest(t, &request, &response)
	if err != nil {
		return "", "", err
	}

	header, err := self.decodeBlockHeader(response.Blocks[0])
	if err != nil {
		return "", "", err
	}

	return response.Blocks[0].HeaderSignature, header.StateRootHash, nil
}

func (self *SawtoothClientTransportZmq) decodeBlockHeader(block *block_pb2.Block) (*block_pb2.BlockHeader, error) {
	var header block_pb2.BlockHeader
	err := proto.Unmarshal(block.Header, &header)
	if err != nil {
		return nil, err
	}

	return &header, nil
}

// GetState returns the state at the given address.
func (self *SawtoothClientTransportZmq) GetState(address string) (*types.State, error) {
	head, stateRoot, err := self.currentStateRoot()
	if err != nil {
		return nil, err
	}

	return self.getStateAtRoot(address, head, stateRoot)
}

// GetStateAtHead returns the state at the given address, at the given head.
func (self *SawtoothClientTransportZmq) GetStateAtHead(address string, head string) (*types.State, error) {
	stateRoot, err := self.headToStateRoot(head)
	if err != nil {
		return nil, err
	}

	return self.getStateAtRoot(address, head, stateRoot)
}

// getStateAtRoot is used to implement both GetState() and GetStateAtHead()
func (self *SawtoothClientTransportZmq) getStateAtRoot(address string, head string, stateRoot string) (*types.State, error) {
	// Set up the request
	t := validator_pb2.Message_CLIENT_STATE_GET_REQUEST
	request := client_state_pb2.ClientStateGetRequest{
		StateRoot: stateRoot,
		Address: address,
	}

	// Send the request and get the response
	var response client_state_pb2.ClientStateGetResponse
	err := self.doZmqRequest(t, &request, &response)
	if err != nil {
		return nil, err
	}

	// Populate our State object
	state := types.State{
		Data: response.Value,
		Address: address,
		Head: head,
	}

	return &state, nil
}

type stateZmqIterator struct {
	commonZmqIterator
	head		string
	stateRoot	string
	address		string
}

func (self *stateZmqIterator) Current() (*types.State, error) {
	err := self.checkCurrent()
	if err != nil {
		return nil, err
	}

	data := self.getCurrent().(*types.State)
	return data, nil
}

func (self *stateZmqIterator) BuildRequest(pagingControl *client_list_control_pb2.ClientPagingControls, sortControl []*client_list_control_pb2.ClientSortControls) (validator_pb2.Message_MessageType, proto.Message, proto.Message) {
	t := validator_pb2.Message_CLIENT_STATE_LIST_REQUEST
	request := client_state_pb2.ClientStateListRequest{Address: self.address, StateRoot: self.stateRoot, Paging: pagingControl, Sorting: sortControl}
	response := client_state_pb2.ClientStateListResponse{}
	return t, &request, &response
}

func (self *stateZmqIterator) ParseProto(message proto.Message) ([]interface{}, error) {
	response := message.(*client_state_pb2.ClientStateListResponse)
	result := make([]interface{}, len(response.Entries))

	for i, item := range response.Entries {
		state := &types.State{
			Data: item.Data,
			Address: item.Address,
			Head: self.head,
		}
		result[i] = state
	}

	return result, nil
}

// GetStateData returns a types.StateIterator that can iterate over all state matching the given prefix.
func (self *SawtoothClientTransportZmq) GetStateIterator(addressPrefix string, fetch int, reverse bool) types.StateIterator {
	pagingControl := &client_list_control_pb2.ClientPagingControls{
		Limit: int32(fetch),
	}

	sortControl := []*client_list_control_pb2.ClientSortControls{
		{
			Keys: []string{"default"},
			Reverse: reverse,
		},
	}

	iterator := &stateZmqIterator{address: addressPrefix}
	iterator.commonZmqIterator = *NewCommonZmqIterator(self, pagingControl, sortControl, iterator)

	head, stateRoot, err := self.currentStateRoot()
	if err != nil {
		iterator.err = err
		return iterator
	}

	iterator.head = head
	iterator.stateRoot = stateRoot

	return iterator
}
