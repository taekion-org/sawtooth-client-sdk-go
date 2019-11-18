package zmq

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/hyperledger/sawtooth-sdk-go/protobuf/validator_pb2"
)

type SawtoothClientTransportZmqError struct {
	// Type is the type of message that was sent to the validator
	Type		validator_pb2.Message_MessageType
	// Request is the actual request sent to the validator
	Request		proto.Message
	// Response is the actual response returned from the validator
	Response	proto.Message
}

// NewSawtoothClientTransportZmqError constructs a SawtoothClientTransportZmqError
func NewSawtoothClientTransportZmqError(t validator_pb2.Message_MessageType, request proto.Message, response proto.Message) *SawtoothClientTransportZmqError {
	return &SawtoothClientTransportZmqError{
		Type:     t,
		Request:  request,
		Response: response,
	}
}

// Error implements the error interface for SawtoothClientTransportZmqError
func (self *SawtoothClientTransportZmqError) Error() string {
	msg := fmt.Sprintf("Sawtooth ZMQ Error: %s", self.Type, self.Request, self.Response)
	return msg
}