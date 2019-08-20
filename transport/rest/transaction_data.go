package rest

import (
	"encoding/json"
	"fmt"
	"github.com/taekion-org/sawtooth-client-sdk-go/transport/types"
	"net/url"
)

// transactionRestResponseSingle represents a REST API reply when a single transaction is requested.
type transactionRestResponseSingle struct {
	Data types.Transaction `json:"data"`
	Link string            `json:"link"`
}

// transactionRestResponseMultiple represents a REST API reply when multiple transactions are requested.
type transactionRestResponseMultiple struct {
	Data	[]types.Transaction `json:"data"`
}

// GetTransaction returns the transaction represented by transactionId.
func (self *SawtoothClientTransportRest) GetTransaction(transactionId string) (*types.Transaction, error) {
	relativeUrl := &url.URL{Path: fmt.Sprintf("/transactions/%s", transactionId)}

	data, err := self.doGetRequest(relativeUrl)
	if err != nil {
		return nil, err
	}

	var response transactionRestResponseSingle
	err = json.Unmarshal(data, &response)
	if err != nil {
		return nil, err
	}

	return &response.Data, nil
}

// batchRestIterator extends commonRestIterator and implements the types.TransactionIterator interface.
type transactionRestIterator struct {
	commonRestIterator
}

// Current returns the "current" transaction from the iterator.
func (self *transactionRestIterator) Current() (*types.Transaction, error) {
	err := self.checkCurrent()
	if err != nil {
		return nil, err
	}

	data := self.getCurrent().(*types.Transaction)
	return data, nil
}

// UnmarshalData handles unmarshaling the raw transaction data returned from the API.
func (self *transactionRestIterator) UnmarshalData(bytes []byte) ([]interface{}, error) {
	var response transactionRestResponseMultiple
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

// GetTransactionIterator returns a types.TransactionIterator that can iterate over all transactions.
func (self *SawtoothClientTransportRest) GetTransactionIterator(fetch int, reverse bool) types.TransactionIterator {
	relativeUrl := &url.URL{Path: "/transactions"}

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

	iterator := &transactionRestIterator{}
	iterator.commonRestIterator = *NewCommonRestIterator(self, relativeUrl, iterator)

	return iterator
}
