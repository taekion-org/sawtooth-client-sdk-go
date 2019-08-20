package rest

import (
	"encoding/json"
	"fmt"
	"github.com/taekion-org/sawtooth-client-sdk-go/transport/types"
	"net/url"
)

// batchRestResponseSingle represents a REST API reply when a single batch is requested.
type batchRestResponseSingle struct {
	Data types.Batch `json:"data"`
	Link string      `json:"link"`
}

// batchRestResponseMultiple represents a REST API reply when multiple batches are requested.
type batchRestResponseMultiple struct {
	Data	[]types.Batch `json:"data"`
}

// GetBatch returns the batch represented by batchId.
func (self *SawtoothClientTransportRest) GetBatch(batchId string) (*types.Batch, error) {
	relativeUrl := &url.URL{Path: fmt.Sprintf("/batches/%s", batchId)}

	data, err := self.doGetRequest(relativeUrl)
	if err != nil {
		return nil, err
	}

	var response batchRestResponseSingle
	err = json.Unmarshal(data, &response)
	if err != nil {
		return nil, err
	}

	return &response.Data, nil
}

// batchRestIterator extends commonRestIterator and implements the types.BatchIterator interface.
type batchRestIterator struct {
	commonRestIterator
}

// Current returns the "current" batch from the iterator.
func (self *batchRestIterator) Current() (*types.Batch, error) {
	err := self.checkCurrent()
	if err != nil {
		return nil, err
	}

	data := self.getCurrent().(*types.Batch)
	return data, nil
}

// UnmarshalData handles unmarshaling the raw batch data returned from the API.
func (self *batchRestIterator) UnmarshalData(bytes []byte) ([]interface{}, error) {
	var response batchRestResponseMultiple
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

// GetBatchIterator returns a types.BatchIterator that can iterate over all batches.
func (self *SawtoothClientTransportRest) GetBatchIterator(fetch int, reverse bool) types.BatchIterator {
	relativeUrl := &url.URL{Path: "/batches"}

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

	iterator := &batchRestIterator{}
	iterator.commonRestIterator = *NewCommonRestIterator(self, relativeUrl, iterator)

	return iterator
}
