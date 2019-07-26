package intkey

import "github.com/taekion-org/sawtooth-client-sdk-go"

// GetAddressPrefix returns the intkey address prefix.
func GetAddressPrefix() string {
	return sawtooth_client_sdk_go.Hexdigest(FAMILY_NAME)[:FAMILY_NAMESPACE_ADDRESS_LENGTH]
}

// GetAddress returns the address of a particular intkey key.
func GetAddress(name string) string {
	prefix := GetAddressPrefix()
	nameAddress := sawtooth_client_sdk_go.Hexdigest(name)[FAMILY_VERB_ADDRESS_LENGTH:]
	return prefix + nameAddress
}
