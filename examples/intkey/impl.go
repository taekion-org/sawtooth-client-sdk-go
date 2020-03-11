package intkey

import cbor "github.com/brianolson/cbor_go"

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
func (self *IntkeyClientImpl) EncodePayload(payload interface{}) ([]byte, error) {
	return cbor.Dumps(payload)
}

// EncodePayload unmarshals a payload from bytes (in this case, using CBOR).
func (self *IntkeyClientImpl) DecodePayload(bytes []byte, ptr interface{}) error {
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

// GetPayloadInputAddresses returns the addresses of any state that the given payload is
// expected to access during processing.
func (self *IntkeyClientImpl) GetPayloadInputAddresses(payload interface{}) []string {
	name := payload.(*IntkeyPayload).Name
	return []string{GetAddress(name)}
}

// GetPayloadOutputAddresses returns the addresses of any state that the given payload is
// expected to modify during processing.
func (self *IntkeyClientImpl) GetPayloadOutputAddresses(payload interface{}) []string {
	name := payload.(*IntkeyPayload).Name
	return []string{GetAddress(name)}
}

// GetPayloadDependencies returns the IDs for any transactions that the given payload
// depends on. For intkey, this is always an empty list.
func (self *IntkeyClientImpl) GetPayloadDependencies(payload interface{}) []string {
	return []string{}
}
