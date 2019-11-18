// Package intkey provides an alternative intkey client implementation using sawtooth-client-sdk-go.
package intkey

import (
	"fmt"
	"github.com/taekion-org/sawtooth-client-sdk-go"
	"github.com/taekion-org/sawtooth-client-sdk-go/transport"
)

// IntkeyPayload represents a payload meant for the intkey transaction processor.
type IntkeyPayload struct {
	Verb		string
	Name		string
	Value		uint
}

// IntKeyClient is the client library for intkey.
type IntkeyClient struct {
	*sawtooth_client_sdk_go.SawtoothClient
}

// NewIntkeyClient returns a new instance of IntkeyClient.
func NewIntkeyClient(url string, keyFile string) (*IntkeyClient, error) {
	args := &sawtooth_client_sdk_go.SawtoothClientArgs{
		URL: url,
		KeyFile: keyFile,
		TransportType: transport.TRANSPORT_REST,
		Impl: &IntkeyClientImpl{},
	}

	sawtoothClient, err := sawtooth_client_sdk_go.NewClient(args)
	if err != nil {
		return nil, err
	}

	client := &IntkeyClient{
		SawtoothClient: sawtoothClient,
	}

	return client, nil
}

// NewIntkeyClient returns a new instance of IntkeyClient that uses the ZMQ transport.
func NewIntkeyClientZmq(url string, keyFile string) (*IntkeyClient, error) {
	args := &sawtooth_client_sdk_go.SawtoothClientArgs{
		URL: url,
		KeyFile: keyFile,
		TransportType: transport.TRANSPORT_ZMQ,
		Impl: &IntkeyClientImpl{},
	}

	sawtoothClient, err := sawtooth_client_sdk_go.NewClient(args)
	if err != nil {
		return nil, err
	}

	client := &IntkeyClient{
		SawtoothClient: sawtoothClient,
	}

	return client, nil
}

// List returns the current mapping of keys to values.
func (self *IntkeyClient) List() (map[string]uint, error) {
	addressPrefix := GetAddressPrefix()
	iterator := self.Transport.GetStateIterator(addressPrefix, 10, false)
	result := make(map[string]uint)

	for iterator.Next() {
		err := iterator.Error()
		if err != nil {
			return nil, err
		}

		current, err := iterator.Current()
		if err != nil {
			return nil, err
		}

		m := make(map[string]uint, 1)
		err = self.ClientImpl.DecodeData(current.Data, &m)
		if err != nil {
			return nil, err
		}

		for key, value := range m {
			result[key] = value
		}
	}

	return result, nil
}

// Show returns the current value of a particular key.
func (self *IntkeyClient) Show(name string) (uint, error) {
	address := GetAddress(name)
	state, err := self.Transport.GetState(address)
	if err != nil {
		return 0, err
	}

	m := make(map[string]uint, 1)
	err = self.ClientImpl.DecodeData(state.Data, &m)
	if err != nil {
		return 0, err
	}

	if result, ok := m[name]; ok {
		return result, nil
	} else {
		return 0, fmt.Errorf("Key not found")
	}
}

// sendTransaction is a common method used to construct a payload and submit it for processing.
func (self *IntkeyClient) sendTransaction(verb string, name string, value uint, wait uint) (string, error) {
	payload := IntkeyPayload{
		Verb: verb,
		Name: name,
		Value: value,
	}

	batchId, err := self.ExecutePayload(&payload)
	if err != nil {
		return batchId, err
	}

	if wait > 0 {
		self.WaitBatch(batchId, int(wait), 1)
	}

	return batchId, nil
}

// Set creates a new key -> value mapping.
func (self *IntkeyClient) Set(name string, value uint, wait uint) (string, error) {
	return self.sendTransaction(VERB_SET, name, value, wait)
}

// Inc increments a key's current value by the given parameter.
func (self *IntkeyClient) Inc(name string, value uint, wait uint) (string, error) {
	return self.sendTransaction(VERB_INC, name, value, wait)
}

// Dec decrements a key's current value by the given parameter.
func (self *IntkeyClient) Dec(name string, value uint, wait uint) (string, error) {
	return self.sendTransaction(VERB_DEC, name, value, wait)
}

// Status returns the status of a given batch.
func (self *IntkeyClient) Status(batchId string) (string, error) {
	status, err := self.Transport.GetBatchStatus(batchId, 0)
	return string(status), err
}
