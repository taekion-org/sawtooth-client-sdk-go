// Package rest provides a SawtoothClientTransport implementation for the Sawtooth REST API.
package rest

import (
	"net/http"
	"net/url"
	"time"
)

// HTTP_TIMEOUT is the default timeout.
const HTTP_TIMEOUT = time.Second * 60

// SawtoothClientTransportRest represents a connection to the REST API.
type SawtoothClientTransportRest struct {
	// URL is the URL to the REST API.
	URL			*url.URL
	// HttpClient is the client that is maintained throughout the life of this object.
	HttpClient	*http.Client
}

// NewSawtoothClientTransportRest returns a new SawtoothClientTransportRest for the given URL.
// Returns an error if a test request to the API does not succeed.
func NewSawtoothClientTransportRest(url *url.URL) (*SawtoothClientTransportRest, error) {
	client := &SawtoothClientTransportRest{
		URL: url,
		HttpClient: &http.Client{Timeout: HTTP_TIMEOUT},
	}

	err := client.testConnection()
	if err != nil {
		return nil, err
	}

	return client, nil
}

// Do the simplest possible request to verify REST API connectivity
func (self *SawtoothClientTransportRest) testConnection() error {
	relativeUrl := &url.URL{Path: "/peers"}

	_, err := self.doGetRequest(relativeUrl)
	if err != nil {
		return err
	}

	return nil
}
