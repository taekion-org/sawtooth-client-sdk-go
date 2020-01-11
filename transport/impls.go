package transport

import (
	"fmt"
	"github.com/taekion-org/sawtooth-client-sdk-go/transport/rest"
	"github.com/taekion-org/sawtooth-client-sdk-go/transport/zmq"
	"net/url"
)

// SawtoothClientTransportType represents an individual transport implementation.
type SawtoothClientTransportType string

// TRANSPORT_REST represents the REST API transport implementation.
const TRANSPORT_REST SawtoothClientTransportType = "rest"
// TRANSPORT_ZMQ represents the ZMQ transport implementation.
const TRANSPORT_ZMQ SawtoothClientTransportType = "zmq"

// NewSawtoothClientTransport instantiates and returns a new SawtoothClientTransport of the specified type.
func NewSawtoothClientTransport(transportType SawtoothClientTransportType, url *url.URL) (SawtoothClientTransport, error) {
	switch transportType {
	case TRANSPORT_REST:
		return rest.NewSawtoothClientTransportRest(url)
	case TRANSPORT_ZMQ:
		return zmq.NewSawtoothClientTransportZmq(url)
	default:
		return nil, fmt.Errorf("Unknown transport type")
	}
}
