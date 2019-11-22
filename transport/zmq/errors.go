package zmq

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/hyperledger/sawtooth-sdk-go/protobuf/client_batch_pb2"
	"github.com/hyperledger/sawtooth-sdk-go/protobuf/client_batch_submit_pb2"
	"github.com/hyperledger/sawtooth-sdk-go/protobuf/client_block_pb2"
	client_peer "github.com/hyperledger/sawtooth-sdk-go/protobuf/client_peers_pb2"
	"github.com/hyperledger/sawtooth-sdk-go/protobuf/client_state_pb2"
	"github.com/hyperledger/sawtooth-sdk-go/protobuf/client_transaction_pb2"
	"github.com/hyperledger/sawtooth-sdk-go/protobuf/validator_pb2"
	"github.com/taekion-org/sawtooth-client-sdk-go/transport/errors"
)

type SawtoothClientTransportZmqError struct {
	// Type is the type of message that was sent to the validator
	Type		validator_pb2.Message_MessageType
	// Request is the actual request sent to the validator
	Request		proto.Message
	// Response is the actual response returned from the validator
	Response	proto.Message
	// ErrorCode is the Sawtooth error code
	ErrorCode	errors.SawtoothTransportErrorCode
}


// Error implements the error interface for SawtoothClientTransportZmqError
func (self *SawtoothClientTransportZmqError) Error() string {
	msg := fmt.Sprintf("Sawtooth ZMQ Error - MessageType: %s - ErrorCode - %d", self.Type, self.ErrorCode)
	return msg
}

// NewSawtoothClientTransportZmqError constructs a SawtoothClientTransportZmqError
func NewSawtoothClientTransportZmqError(t validator_pb2.Message_MessageType, request proto.Message, response proto.Message, errorCode errors.SawtoothTransportErrorCode) *SawtoothClientTransportZmqError {
	return &SawtoothClientTransportZmqError{
		Type: t,
		Request: request,
		Response: response,
		ErrorCode: errorCode,
	}
}

func checkForError(response proto.Message) errors.SawtoothTransportErrorCode {
	switch r := response.(type) {
	case *client_batch_pb2.ClientBatchGetResponse:
		switch r.Status {
		case client_batch_pb2.ClientBatchGetResponse_OK:
			return errors.NO_ERROR
		case client_batch_pb2.ClientBatchGetResponse_INTERNAL_ERROR:
			return errors.VALIDATOR_UNKNOWN_ERROR
		case client_batch_pb2.ClientBatchGetResponse_INVALID_ID:
			return errors.INVALID_RESOURCE_ID
		case client_batch_pb2.ClientBatchGetResponse_NO_RESOURCE:
			return errors.BATCH_NOT_FOUND
		}

	case *client_batch_pb2.ClientBatchListResponse:
		switch r.Status {
		case client_batch_pb2.ClientBatchListResponse_OK:
			return errors.NO_ERROR
		case client_batch_pb2.ClientBatchListResponse_INTERNAL_ERROR:
			return errors.VALIDATOR_UNKNOWN_ERROR
		case client_batch_pb2.ClientBatchListResponse_NOT_READY:
			return errors.VALIDATOR_NOT_READY
		case client_batch_pb2.ClientBatchListResponse_NO_ROOT:
			return errors.INVALID_HEAD
		case client_batch_pb2.ClientBatchListResponse_INVALID_PAGING:
			return errors.INVALID_PAGING_QUERY
		case client_batch_pb2.ClientBatchListResponse_INVALID_SORT:
			return errors.INVALID_SORT_QUERY
		case client_batch_pb2.ClientBatchListResponse_INVALID_ID:
			return errors.INVALID_RESOURCE_ID
		case client_batch_pb2.ClientBatchListResponse_NO_RESOURCE:
			return errors.BATCH_NOT_FOUND
		}

	case *client_batch_submit_pb2.ClientBatchStatusResponse:
		switch r.Status {
		case client_batch_submit_pb2.ClientBatchStatusResponse_OK:
			return errors.NO_ERROR
		case client_batch_submit_pb2.ClientBatchStatusResponse_INTERNAL_ERROR:
			return errors.VALIDATOR_UNKNOWN_ERROR
		case client_batch_submit_pb2.ClientBatchStatusResponse_INVALID_ID:
			return errors.INVALID_RESOURCE_ID
		case client_batch_submit_pb2.ClientBatchStatusResponse_NO_RESOURCE:
			return errors.BATCH_STATUS_UNAVAILABLE
		}

	case *client_batch_submit_pb2.ClientBatchSubmitResponse:
		switch r.Status {
		case client_batch_submit_pb2.ClientBatchSubmitResponse_OK:
			return errors.NO_ERROR
		case client_batch_submit_pb2.ClientBatchSubmitResponse_INTERNAL_ERROR:
			return errors.VALIDATOR_UNKNOWN_ERROR
		case client_batch_submit_pb2.ClientBatchSubmitResponse_INVALID_BATCH:
			return errors.BATCH_INVALID
		case client_batch_submit_pb2.ClientBatchSubmitResponse_QUEUE_FULL:
			return errors.BATCH_UNABLE_TO_ACCEPT
		}

	case *client_block_pb2.ClientBlockGetResponse:
		switch r.Status {
		case client_block_pb2.ClientBlockGetResponse_OK:
			return errors.NO_ERROR
		case client_block_pb2.ClientBlockGetResponse_INTERNAL_ERROR:
			return errors.VALIDATOR_UNKNOWN_ERROR
		case client_block_pb2.ClientBlockGetResponse_INVALID_ID:
			return errors.INVALID_RESOURCE_ID
		case client_block_pb2.ClientBlockGetResponse_NO_RESOURCE:
			return errors.BLOCK_NOT_FOUND
		}

	case *client_block_pb2.ClientBlockListResponse:
		switch r.Status {
		case client_block_pb2.ClientBlockListResponse_OK:
			return errors.NO_ERROR
		case client_block_pb2.ClientBlockListResponse_INTERNAL_ERROR:
			return errors.VALIDATOR_UNKNOWN_ERROR
		case client_block_pb2.ClientBlockListResponse_NOT_READY:
			return errors.VALIDATOR_NOT_READY
		case client_block_pb2.ClientBlockListResponse_NO_ROOT:
			return errors.INVALID_HEAD
		case client_block_pb2.ClientBlockListResponse_INVALID_PAGING:
			return errors.INVALID_PAGING_QUERY
		case client_block_pb2.ClientBlockListResponse_INVALID_SORT:
			return errors.INVALID_SORT_QUERY
		case client_block_pb2.ClientBlockListResponse_INVALID_ID:
			return errors.INVALID_RESOURCE_ID
		case client_block_pb2.ClientBlockListResponse_NO_RESOURCE:
			return errors.BLOCK_NOT_FOUND
		}

	case *client_transaction_pb2.ClientTransactionGetResponse:
		switch r.Status {
		case client_transaction_pb2.ClientTransactionGetResponse_OK:
			return errors.NO_ERROR
		case client_transaction_pb2.ClientTransactionGetResponse_INTERNAL_ERROR:
			return errors.VALIDATOR_UNKNOWN_ERROR
		case client_transaction_pb2.ClientTransactionGetResponse_INVALID_ID:
			return errors.INVALID_RESOURCE_ID
		case client_transaction_pb2.ClientTransactionGetResponse_NO_RESOURCE:
			return errors.TRANSACTION_NOT_FOUND
		}

	case *client_transaction_pb2.ClientTransactionListResponse:
		switch r.Status {
		case client_transaction_pb2.ClientTransactionListResponse_OK:
			return errors.NO_ERROR
		case client_transaction_pb2.ClientTransactionListResponse_INTERNAL_ERROR:
			return errors.VALIDATOR_UNKNOWN_ERROR
		case client_transaction_pb2.ClientTransactionListResponse_NOT_READY:
			return errors.VALIDATOR_NOT_READY
		case client_transaction_pb2.ClientTransactionListResponse_NO_ROOT:
			return errors.INVALID_HEAD
		case client_transaction_pb2.ClientTransactionListResponse_INVALID_PAGING:
			return errors.INVALID_PAGING_QUERY
		case client_transaction_pb2.ClientTransactionListResponse_INVALID_SORT:
			return errors.INVALID_SORT_QUERY
		case client_transaction_pb2.ClientTransactionListResponse_INVALID_ID:
			return errors.INVALID_RESOURCE_ID
		case client_transaction_pb2.ClientTransactionListResponse_NO_RESOURCE:
			return errors.TRANSACTION_NOT_FOUND
		}

	case *client_state_pb2.ClientStateGetResponse:
		switch r.Status {
		case client_state_pb2.ClientStateGetResponse_OK:
			return errors.NO_ERROR
		case client_state_pb2.ClientStateGetResponse_INTERNAL_ERROR:
			return errors.VALIDATOR_UNKNOWN_ERROR
		case client_state_pb2.ClientStateGetResponse_NOT_READY:
			return errors.VALIDATOR_NOT_READY
		case client_state_pb2.ClientStateGetResponse_NO_ROOT, client_state_pb2.ClientStateGetResponse_INVALID_ROOT:
			return errors.INVALID_HEAD
		case client_state_pb2.ClientStateGetResponse_INVALID_ADDRESS:
			return errors.INVALID_STATE_ADDRESS
		case client_state_pb2.ClientStateGetResponse_NO_RESOURCE:
			return errors.STATE_NOT_FOUND
		}

	case *client_state_pb2.ClientStateListResponse:
		switch r.Status {
		case client_state_pb2.ClientStateListResponse_OK:
			return errors.NO_ERROR
		case client_state_pb2.ClientStateListResponse_INTERNAL_ERROR:
			return errors.VALIDATOR_UNKNOWN_ERROR
		case client_state_pb2.ClientStateListResponse_NOT_READY:
			return errors.VALIDATOR_NOT_READY
		case client_state_pb2.ClientStateListResponse_NO_ROOT, client_state_pb2.ClientStateListResponse_INVALID_ROOT:
			return errors.INVALID_HEAD
		case client_state_pb2.ClientStateListResponse_INVALID_PAGING:
			return errors.INVALID_PAGING_QUERY
		case client_state_pb2.ClientStateListResponse_INVALID_SORT:
			return errors.INVALID_SORT_QUERY
		case client_state_pb2.ClientStateListResponse_INVALID_ADDRESS:
			return errors.INVALID_STATE_ADDRESS
		case client_state_pb2.ClientStateListResponse_NO_RESOURCE:
			return errors.STATE_NOT_FOUND
		}
	case *client_peer.ClientPeersGetResponse:
		switch r.Status {
		case client_peer.ClientPeersGetResponse_OK:
			return errors.NO_ERROR
		case client_peer.ClientPeersGetResponse_ERROR:
			return errors.VALIDATOR_UNKNOWN_ERROR
		}

	}

	return errors.UNKNOWN_ERROR
}
