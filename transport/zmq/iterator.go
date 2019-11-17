package zmq

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/hyperledger/sawtooth-sdk-go/protobuf/client_list_control_pb2"
	"github.com/hyperledger/sawtooth-sdk-go/protobuf/validator_pb2"
)

// zmqIteratorImpl must be implemented by all ZMQ iterators.
type zmqIteratorImpl interface {
	BuildRequest(*client_list_control_pb2.ClientPagingControls, []*client_list_control_pb2.ClientSortControls) (validator_pb2.Message_MessageType, proto.Message, proto.Message)
	ParseProto(message proto.Message) ([]interface{}, error)
}

// zmqPagingResponseGetter is a special interface used to extract the paging info
// from different proto buffer objects.
type zmqPagingResponseGetter interface {
	GetPaging() *client_list_control_pb2.ClientPagingResponse
}

// commonZmqIterator implements an iterator for the validator ZMQ interface that can be extended to be used
// across multiple object types.
type commonZmqIterator struct {
	transport			*SawtoothClientTransportZmq
	nextPagingControl	*client_list_control_pb2.ClientPagingControls
	sortControl			[]*client_list_control_pb2.ClientSortControls

	data		[]interface{}
	current		interface{}
	err			error

	impl		zmqIteratorImpl
}

// NewCommonZmqIterator returns a new commonZmqIterator for use in composing a usable object iterator.
func NewCommonZmqIterator(transport *SawtoothClientTransportZmq,
							pagingControl *client_list_control_pb2.ClientPagingControls,
							sortControl []*client_list_control_pb2.ClientSortControls,
							impl zmqIteratorImpl) *commonZmqIterator {
	return &commonZmqIterator{transport: transport, nextPagingControl: pagingControl, sortControl: sortControl, impl: impl}
}

// Next returns true if a next value is available.
func (self *commonZmqIterator) Next() bool {
	err := self.fetchNext()
	if err != nil {
		self.err = err
		return false
	}

	if len(self.data) == 0 {
		return false
	}

	// Pop and shift
	self.current, self.data = self.data[0], self.data[1:]

	return true
}

// Error returns the error (if any) contained in the iterator.
func (self *commonZmqIterator) Error() error {
	return self.err
}

// fetchNext get the next value or batch of values from the validator.
func (self *commonZmqIterator) fetchNext() error {
	if len(self.data) > 0 {
		return nil
	}

	if self.nextPagingControl == nil {
		return nil
	}

	if self.err != nil {
		return nil
	}

	// Do the ZMQ request
	t, requestMsg, responseMsg := self.impl.BuildRequest(self.nextPagingControl, self.sortControl)
	err := self.transport.doZmqRequest(t, requestMsg, responseMsg)
	if err != nil {
		return nil
	}

	// Parse out the actual data
	self.data, err = self.impl.ParseProto(responseMsg)
	if err != nil {
		return err
	}

	// Get the next paging info
	pagingResponse := responseMsg.(zmqPagingResponseGetter).GetPaging()
	if pagingResponse.Next == "" {
		self.nextPagingControl = nil
	} else {
		self.nextPagingControl = &client_list_control_pb2.ClientPagingControls{
			Start: pagingResponse.Next,
			Limit: self.nextPagingControl.Limit,
		}
	}

	return nil
}

// checkCurrent checks to make sure there is a current value in the iterator. If no current
// value is present, returns an error.
func (self *commonZmqIterator) checkCurrent() error {
	if self.current == nil {
		return fmt.Errorf("No current value in iterator...")
	}

	return nil
}

// getCurrent returns the current value from the iterator as an interface{} type.
func (self *commonZmqIterator) getCurrent() interface{} {
	return self.current
}
