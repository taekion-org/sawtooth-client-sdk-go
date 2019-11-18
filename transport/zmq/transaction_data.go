package zmq

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/hyperledger/sawtooth-sdk-go/protobuf/client_list_control_pb2"
	"github.com/hyperledger/sawtooth-sdk-go/protobuf/client_transaction_pb2"
	"github.com/hyperledger/sawtooth-sdk-go/protobuf/validator_pb2"
	"github.com/taekion-org/sawtooth-client-sdk-go/transport/types"
)

func (self *SawtoothClientTransportZmq) GetTransaction(transactionId string) (*types.Transaction, error) {
	// Set up the request
	t := validator_pb2.Message_CLIENT_TRANSACTION_GET_REQUEST
	request := client_transaction_pb2.ClientTransactionGetRequest{
		TransactionId: transactionId,
	}

	// Send the request and get the response
	var response client_transaction_pb2.ClientTransactionGetResponse
	err := self.doZmqRequest(t, &request, &response)
	if err != nil {
		return nil, err
	}

	// Convert the resulting batch into our data type
	transaction, err := types.TransactionFromProto(response.Transaction)
	if err != nil {
		return nil, fmt.Errorf("Error parsing transaction protobuf: %s", err)
	}

	// Return the result
	return transaction, nil
}

type transactionZmqIterator struct {
	commonZmqIterator
}

func (self *transactionZmqIterator) Current() (*types.Transaction, error) {
	err := self.checkCurrent()
	if err != nil {
		return nil, err
	}

	data := self.getCurrent().(*types.Transaction)
	return data, nil
}

func (self *transactionZmqIterator) BuildRequest(pagingControl *client_list_control_pb2.ClientPagingControls, sortControl []*client_list_control_pb2.ClientSortControls) (validator_pb2.Message_MessageType, proto.Message, proto.Message) {
	t := validator_pb2.Message_CLIENT_TRANSACTION_LIST_REQUEST
	request := client_transaction_pb2.ClientTransactionListRequest{Paging: pagingControl, Sorting: sortControl}
	response := client_transaction_pb2.ClientTransactionListResponse{}
	return t, &request, &response
}

func (self *transactionZmqIterator) ParseProto(message proto.Message) ([]interface{}, error) {
	response := message.(*client_transaction_pb2.ClientTransactionListResponse)
	result := make([]interface{}, len(response.Transactions))

	for i, item := range response.Transactions {
		transaction, err := types.TransactionFromProto(item)
		if err != nil {
			return nil, err
		}
		result[i] = transaction
	}

	return result, nil
}

// GetTransactionIterator returns a types.TransactionIterator that can iterate over all transactions.
func (self *SawtoothClientTransportZmq) GetTransactionIterator(fetch int, reverse bool) types.TransactionIterator {
	pagingControl := &client_list_control_pb2.ClientPagingControls{
		Limit: int32(fetch),
	}

	sortControl := []*client_list_control_pb2.ClientSortControls{
		{
			Keys: []string{"default"},
			Reverse: reverse,
		},
	}

	iterator := &transactionZmqIterator{}
	iterator.commonZmqIterator = *NewCommonZmqIterator(self, pagingControl, sortControl, iterator)

	return iterator
}
