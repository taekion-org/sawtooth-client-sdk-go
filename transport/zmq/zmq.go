// Package zmq provides a SawtoothClientTransport implementation for the Sawtooth validator ZMQ interface.
package zmq

import (
	"fmt"
	"net/url"
	"github.com/pebbe/zmq4"
	"github.com/hyperledger/sawtooth-sdk-go/messaging"
	"github.com/hyperledger/sawtooth-sdk-go/protobuf/client_peers_pb2"
	"github.com/hyperledger/sawtooth-sdk-go/protobuf/validator_pb2"
	log "github.com/sirupsen/logrus"
)

const MAX_CONNECTIONS = 8

type sawtoothZmqConnection struct {
	messaging.Connection
}

func (self *sawtoothZmqConnection) String() string {
	return fmt.Sprintf("sawtoothZmqConnection(identity=%s)", self.Connection.Identity())
}

// SawtoothClientTransportZmq represents a connection to the validator via ZMQ.
type SawtoothClientTransportZmq struct {
	// URL is the ZMQ URL to the validator
	URL			*url.URL

	// Context is a common ZMQ context for the transport
	Context		*zmq4.Context

	// ConnectionChannel is a channel to hold ZMQ connections
	ConnectionChannel		chan *sawtoothZmqConnection

	// Context for logging
	logContext	*log.Entry
}

// NewSawtoothClientTransportZmq returns a new SawtoothClientTransportZmq for the given URL.
// Returns an error if a test request to the validator does not succeed.
func NewSawtoothClientTransportZmq(url *url.URL) (*SawtoothClientTransportZmq, error) {
	// Create a new transport object
	client := &SawtoothClientTransportZmq{
		URL: url,
		logContext: log.WithField("object", "SawtoothClientTransportZmq"),
	}

	// Create a ZMQ context
	context, err := zmq4.NewContext()
	if err != nil {
		return nil, err
	}
	client.Context = context

	// Initialize the channel
	client.ConnectionChannel = make(chan *sawtoothZmqConnection, MAX_CONNECTIONS)

	// Test the connection
	err = client.testConnection()
	if err != nil {
		return nil, err
	}

	return client, nil
}

// Create a new ZMQ connection
func (self *SawtoothClientTransportZmq) newConnection() (*sawtoothZmqConnection, error) {
	logContext := self.logContext.WithField("method", "newConnection")
	logContext.Trace()

	rawConn, err := messaging.NewConnection(self.Context, zmq4.DEALER, self.URL.String(), false)
	if err != nil {
		return nil, err
	}

	conn := &sawtoothZmqConnection{
		Connection: rawConn,
	}

	return conn, nil
}

// Get a ZMQ connection from the pool
func (self *SawtoothClientTransportZmq) getConnection() (*sawtoothZmqConnection, error) {
	logContext := self.logContext.WithField("method", "getConnection")
	logContext.Trace()

	var conn *sawtoothZmqConnection

	select {
	case conn = <- self.ConnectionChannel:
		return conn, nil
	default:
		return self.newConnection()
	}
}

// Return a ZMQ connection to the pool
func (self *SawtoothClientTransportZmq) putConnection(conn *sawtoothZmqConnection) error {
	logContext := self.logContext.WithField("method", "putConnection")
	logContext.Trace()

	select {
	case self.ConnectionChannel <- conn:
	default:
		logContext.Tracef("No room in connection pool, closing: %s", conn)
		conn.Close()
	}
	return nil
}

// Do a simple request to verify ZMQ connectivity.
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
