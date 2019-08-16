package rest

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/taekion-org/sawtooth-client-sdk-go/transport/types"
	"net/url"
)

// stateRestResponseSingle represents a REST API reply when a single state item is requested.
type stateRestResponseSingle struct {
	Data	string		`json:"data"`
	Head	string		`json:"head"`
	Link	string		`json:"link"`
}

// stateRestResponseMultiple represents a REST API reply when a state prefix is requested.
type stateRestResponseMultiple struct {
	Data   []struct {
		Data    string		`json:"data"`
		Address string		`json:"address"`
	}	`json:"data"`

	Head string		`json:"head"`
	Link string		`json:"link"`
}

// GetState returns the state at the given address.
func (self *SawtoothClientTransportRest) GetState(address string) (*types.State, error) {
	relativeUrl := &url.URL{Path: fmt.Sprintf("/state/%s", address)}

	data, err := self.doGetRequest(relativeUrl)
	if err != nil {
		return nil, err
	}

	var response stateRestResponseSingle
	err = json.Unmarshal(data, &response)
	if err != nil {
		return nil, err
	}

	dataBytes, err := base64.StdEncoding.DecodeString(response.Data)
	if err != nil {
		return nil, err
	}

	return &types.State{Data: dataBytes, Address: address, Head: response.Head}, nil
}

// stateRestIterator extends commonRestIterator and implements the types.StateIterator interface.
type stateRestIterator struct {
	commonRestIterator
}

// Current returns the "current" state from the iterator.
func (self *stateRestIterator) Current() (*types.State, error) {
	err := self.checkCurrent()
	if err != nil {
		return nil, err
	}

	data := self.getCurrent().(*types.State)
	return data, nil
}

// UnmarshalData handles unmarshaling the raw state data returned from the API.
func (self *stateRestIterator) UnmarshalData(bytes []byte) ([]interface{}, error) {
	var response stateRestResponseMultiple
	err := json.Unmarshal(bytes, &response)
	if err != nil {
		return nil, err
	}

	result := make([]interface{}, len(response.Data))
	for i, item := range response.Data {
		dataBytes, err := base64.StdEncoding.DecodeString(item.Data)
		if err != nil {
			return nil, err
		}
		result[i] = &types.State{Data: dataBytes, Address: item.Address, Head: response.Head}
	}

	return result, nil
}

// GetStateData returns a types.StateIterator that can iterate over all state matching the given prefix.
func (self *SawtoothClientTransportRest) GetStateIterator(addressPrefix string, fetch int, reverse bool) types.StateIterator {
	relativeUrl := &url.URL{Path: "/state"}

	query := relativeUrl.Query()
	query.Add("address", addressPrefix)
	if fetch != 0 {
		query.Add("limit", fmt.Sprintf("%d", fetch))
	}
	if reverse {
		query.Add("reverse", "")
	} else {
		query.Add("reverse", "false")
	}
	relativeUrl.RawQuery = query.Encode()

	iterator := &stateRestIterator{}
	iterator.commonRestIterator = *NewCommonRestIterator(self, relativeUrl, iterator)

	return iterator
}
