package rest

import (
	"encoding/json"
	"fmt"
	"net/url"
)

// commonRestPagingData represents the paging data in a REST API reply.
type commonRestPagingData struct {
	Paging struct {
		Limit        int    `json:"limit"`
		Next         string `json:"next"`
		NextPosition string `json:"next_position"`
		Start        string `json:"start"`
	} `json:"paging"`
}

// restIteratorImpl must be implemented by all rest iterators.
type restIteratorImpl interface {
	UnmarshalData(bytes []byte) ([]interface{}, error)
}

// commonRestIterator implements an iterator for the REST API that can be extended to be used
// across multiple object types.
type commonRestIterator struct {
	transport	*SawtoothClientTransportRest
	nextUrl		*url.URL

	data		[]interface{}
	current		interface{}
	err			error

	impl restIteratorImpl
}

// NewCommonRestIterator returns a new commonRestIterator for use in composing a usable object iterator.
func NewCommonRestIterator(transport *SawtoothClientTransportRest, nextUrl *url.URL, impl restIteratorImpl) *commonRestIterator {
	return &commonRestIterator{transport: transport, nextUrl: nextUrl, impl: impl}
}

// Next returns true if a next value is available.
func (self *commonRestIterator) Next() bool {
	err := self.fetchNext()
	if err != nil {
		self.err = err
		return false
	}

	if len(self.data) == 0 {
		return false
	}

	// Pop and shift
	self.current, self.data = self.data[0], self.data[1:]

	return true
}

// Error returns the error (if any) contained in the iterator.
func (self *commonRestIterator) Error() error {
	return self.err
}

// fetchNext get the next value or batch of values from the REST API.
func (self *commonRestIterator) fetchNext() error {
	if len(self.data) > 0 {
		return nil
	}

	if self.nextUrl == nil {
		return nil
	}

	if self.err != nil {
		return nil
	}

	// Do the request to the api
	bytes, err := self.transport.doGetRequest(self.nextUrl)
	if err != nil {
		return err
	}

	// Unmarshal the actual data
	self.data, err = self.impl.UnmarshalData(bytes)
	if err != nil {
		return err
	}

	// Unmarshal the paging info
	var pagingData commonRestPagingData
	err = json.Unmarshal(bytes, &pagingData)
	nextRawUrl := pagingData.Paging.Next
	if nextRawUrl == "" {
		self.nextUrl = nil
		return nil
	}

	// Parse the next url
	nextUrl, err := url.Parse(nextRawUrl)
	if err != nil {
		self.nextUrl = nil
	} else {
		self.nextUrl = nextUrl
	}

	return nil
}

// checkCurrent checks to make sure there is a current value in the iterator. If no current
// value is present, returns an error.
func (self *commonRestIterator) checkCurrent() error {
	if self.current == nil {
		return fmt.Errorf("No current value in iterator...")
	}

	return nil
}

// getCurrent returns the current value from the iterator as an interface{} type.
func (self *commonRestIterator) getCurrent() interface{} {
	return self.current
}
