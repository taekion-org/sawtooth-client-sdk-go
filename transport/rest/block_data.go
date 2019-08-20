package rest

import (
	"encoding/json"
	"fmt"
	"github.com/taekion-org/sawtooth-client-sdk-go/transport/types"
	"net/url"
)

// blockRestResponseSingle represents a REST API reply when a single block is requested.
type blockRestResponseSingle struct {
	Data types.Block `json:"data"`
	Link string      `json:"link"`
}

// blockRestMultipleResponse represents a REST API reply when multiple blocks are requested.
type blockRestMultipleResponse struct {
	Data	[]types.Block `json:"data"`
}

// GetBlock returns a the block represented by blockId.
func (self *SawtoothClientTransportRest) GetBlock(blockId string) (*types.Block, error) {
	relativeUrl := &url.URL{Path: fmt.Sprintf("/blocks/%s", blockId)}

	data, err := self.doGetRequest(relativeUrl)
	if err != nil {
		return nil, err
	}

	var jsonData blockRestResponseSingle
	err = json.Unmarshal(data, &jsonData)
	if err != nil {
		return nil, err
	}

	return &jsonData.Data, nil
}

// blockRestIterator extends commonRestIterator and implements the types.BlockIterator interface.
type blockRestIterator struct {
	commonRestIterator
}

// Current returns the "current" block from the iterator.
func (self *blockRestIterator) Current() (*types.Block, error) {
	err := self.checkCurrent()
	if err != nil {
		return nil, err
	}

	data := self.getCurrent().(*types.Block)
	return data, nil
}

// UnmarshalData handles unmarshaling the raw block data returned from the API.
func (self *blockRestIterator) UnmarshalData(bytes []byte) ([]interface{}, error) {
	var response blockRestMultipleResponse
	err := json.Unmarshal(bytes, &response)
	if err != nil {
		return nil, err
	}

	result := make([]interface{}, len(response.Data))
	for i, item := range response.Data {
		result[i] = &item
	}

	return result, nil
}

// GetBlockIterator returns a types.BlockIterator that can iterate over all blocks.
func (self *SawtoothClientTransportRest) GetBlockIterator(fetch int, reverse bool) types.BlockIterator {
	relativeUrl := &url.URL{Path: "/blocks"}

	query := relativeUrl.Query()
	if fetch != 0 {
		query.Add("limit", fmt.Sprintf("%d", fetch))
	}
	if reverse {
		query.Add("reverse", "")
	} else {
		query.Add("reverse", "false")
	}
	relativeUrl.RawQuery = query.Encode()

	iterator := &blockRestIterator{}
	iterator.commonRestIterator = *NewCommonRestIterator(self, relativeUrl, iterator)

	return iterator
}
