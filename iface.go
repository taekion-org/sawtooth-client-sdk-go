package sawtooth_client_sdk_go

// SawtoothClientImpl is the interface that must be implemented by each user of the
// library.
type SawtoothClientImpl interface {
	GetFamilyName() string
	GetFamilyVersion() string

	EncodePayload(interface{}) ([]byte, error)
	DecodePayload([]byte, interface{}) error

	EncodeData(interface{}) ([]byte, error)
	DecodeData([]byte, interface{}) error

	GetPayloadInputAddresses(payload interface{}) []string
	GetPayloadOutputAddresses(payload interface{}) []string
}
