// Package error provides a generalized error interface for transport implementations.
package errors

import (
	"fmt"
)

// SawtoothTransportErrorCode represents an error returned by Sawtooth.
type SawtoothTransportErrorCode uint

const (
	NO_ERROR						SawtoothTransportErrorCode		= 0
	UNKNOWN_ERROR					SawtoothTransportErrorCode		= 1024

	VALIDATOR_UNKNOWN_ERROR			SawtoothTransportErrorCode		= 10
	VALIDATOR_NOT_READY				SawtoothTransportErrorCode		= 15
	VALIDATOR_TIMED_OUT				SawtoothTransportErrorCode		= 17
	VALIDATOR_DISCONNECTED			SawtoothTransportErrorCode		= 18
	VALIDATOR_INVALID_RESPONSE		SawtoothTransportErrorCode		= 20

	BATCH_STATUS_UNAVAILABLE		SawtoothTransportErrorCode		= 27
	BATCH_INVALID					SawtoothTransportErrorCode		= 30
	BATCH_UNABLE_TO_ACCEPT			SawtoothTransportErrorCode		= 31
	BATCH_NONE_SUBMITTED			SawtoothTransportErrorCode		= 34
	BATCH_PROTOBUF_NOT_DECODABLE	SawtoothTransportErrorCode		= 35

	INVALID_HEAD					SawtoothTransportErrorCode		= 50
	INVALID_COUNT_QUERY				SawtoothTransportErrorCode		= 53
	INVALID_PAGING_QUERY			SawtoothTransportErrorCode		= 54
	INVALID_SORT_QUERY				SawtoothTransportErrorCode		= 57
	INVALID_RESOURCE_ID				SawtoothTransportErrorCode		= 60
	INVALID_STATE_ADDRESS			SawtoothTransportErrorCode		= 62

	BLOCK_NOT_FOUND					SawtoothTransportErrorCode		= 70
	BATCH_NOT_FOUND					SawtoothTransportErrorCode		= 71
	TRANSACTION_NOT_FOUND			SawtoothTransportErrorCode		= 72
	STATE_NOT_FOUND					SawtoothTransportErrorCode		= 75
	TRANSACTION_RECEIPT_NOT_FOUND	SawtoothTransportErrorCode		= 80
)

// SawtoothClientTransportError represents an error returned by a SawtoothClientTransport
// implementation.
type SawtoothClientTransportError struct {
	ErrorCode      SawtoothTransportErrorCode
	TransportError error
}

// Error implements the error interface for SawtoothClientTransportError.
func (self *SawtoothClientTransportError) Error() string {
	return fmt.Sprintf("Sawtooth Error: %d -- Transport Error: %s", self.ErrorCode, self.TransportError)
}
