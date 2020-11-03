package intkey

import (
	cbor "github.com/brianolson/cbor_go"
	sawtooth_client_sdk_go "github.com/taekion-org/sawtooth-client-sdk-go"
)

// IntkeyClientImpl is the type that implements the required SawtoothClientImpl interface.
type IntkeyClientImpl struct {}

// GetFamilyName returns the family name.
func (self *IntkeyClientImpl) GetFamilyName() string {
	return FAMILY_NAME
}

// GetFamilyVersion returns the family version.
func (self *IntkeyClientImpl) GetFamilyVersion() string {
	return FAMILY_VERSION
}

// EncodePayload marshals a payload into bytes (in this case, using CBOR).
func (self *IntkeyClientImpl) EncodePayload(payload sawtooth_client_sdk_go.SawtoothPayload) ([]byte, error) {
	return cbor.Dumps(payload)
}

// EncodePayload unmarshals a payload from bytes (in this case, using CBOR).
func (self *IntkeyClientImpl) DecodePayload(bytes []byte, ptr sawtooth_client_sdk_go.SawtoothPayload) error {
	return cbor.Loads(bytes, ptr)
}

// EncodeData marshals a data item into bytes (in this case, using CBOR).
func (self *IntkeyClientImpl) EncodeData(data interface{}) ([]byte, error) {
	return cbor.Dumps(data)
}
// EncodePayDecodeDataload unmarshals a data item from bytes (in this case, using CBOR).
func (self *IntkeyClientImpl) DecodeData(bytes []byte, ptr interface{}) error {
	return cbor.Loads(bytes, ptr)
}
