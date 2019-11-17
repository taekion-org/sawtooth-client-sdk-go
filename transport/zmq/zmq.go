// Package zmq provides a SawtoothClientTransport implementation for the Sawtooth validator ZMQ interface.
package zmq

import (
	"net/url"
	"github.com/pebbe/zmq4"
	"github.com/hyperledger/sawtooth-sdk-go/messaging"
	"github.com/hyperledger/sawtooth-sdk-go/protobuf/client_peers_pb2"
	"github.com/hyperledger/sawtooth-sdk-go/protobuf/validator_pb2"
)

// SawtoothClientTransportZmq represents a connection to the validator via ZMQ.
type SawtoothClientTransportZmq struct {
	// URL is the ZMQ URL to the validator
	URL			*url.URL

	// Connection is a ZMQ connection
	Connection	messaging.Connection
}

// NewSawtoothClientTransportZmq returns a new SawtoothClientTransportZmq for the given URL.
// Returns an error if a test request to the validator does not succeed.
func NewSawtoothClientTransportZmq(url *url.URL) (*SawtoothClientTransportZmq, error) {
	client := &SawtoothClientTransportZmq{
		URL: url,
	}

	context, err := zmq4.NewContext()
	if err != nil {
		return nil, err
	}

	// Create a zmq connection
	connection, err := messaging.NewConnection(context, zmq4.DEALER, url.String(), false)
	if err != nil {
		return nil, err
	}
	client.Connection = connection

	err = client.testConnection()
	if err != nil {
		return nil, err
	}

	return client, nil
}

// Do the simplest possible request to verify ZMQ connectivity.
func (self *SawtoothClientTransportZmq) testConnection() error {
	t := validator_pb2.Message_CLIENT_PEERS_GET_REQUEST
	request := client_peer.ClientPeersGetRequest{}
	var response client_peer.ClientPeersGetResponse

	err := self.doZmqRequest(t, &request, &response)
	if err != nil {
		return err
	}

	return nil
}
