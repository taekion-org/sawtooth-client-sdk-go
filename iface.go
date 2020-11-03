package sawtooth_client_sdk_go

// SawtoothClientImpl is the interface that must be implemented by each user of the
// library.
type SawtoothClientImpl interface {
	GetFamilyName() string
	GetFamilyVersion() string

	EncodePayload(SawtoothPayload) ([]byte, error)
	DecodePayload([]byte, SawtoothPayload) error

	EncodeData(interface{}) ([]byte, error)
	DecodeData([]byte, interface{}) error
}

// SawtoothPayload is the interface that must be implemented by payload objects.
type SawtoothPayload interface {
	GetInputAddresses() []string
	GetOutputAddresses() []string
	GetDependencies() []string
	GetNonce() string
}
