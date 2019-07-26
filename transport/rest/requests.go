package rest

import (
	"bytes"
	"github.com/taekion-org/sawtooth-client-sdk-go/transport/errors"
	"io/ioutil"
	"net/http"
	"net/url"
)

// doGetRequest provides a generalized GET call to the REST API. Returns the response as
// a []byte slice, or an error if something goes wrong.
func (self *SawtoothClientTransportRest) doGetRequest(relativeUrl *url.URL) ([]byte, error) {
	fullUrl := self.URL.ResolveReference(relativeUrl)

	request, err := http.NewRequest(http.MethodGet, fullUrl.String(), nil)
	if err != nil {
		return nil, err
	}
	request.Header.Set("Accept", "application/json")

	response, err := self.HttpClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		transportError := NewSawtoothClientTransportRestError(response)
		err = &errors.SawtoothClientTransportError{
			ErrorCode:      errors.SawtoothTransportErrorCode(transportError.ErrorResponse.Error.Code),
			TransportError: transportError,
		}

		return nil, err
	}

	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	return responseData, nil
}

// doPostRequestBinary calls doPostRequest while setting contentType to "application/octet-stream".
func (self *SawtoothClientTransportRest) doPostRequestBinary(relativeUrl *url.URL, data []byte) ([]byte, error) {
	return self.doPostRequest(relativeUrl, data, "application/octet-stream")
}

// doPostRequestBinary calls doPostRequest while setting contentType to "application/json".
func (self *SawtoothClientTransportRest) doPostRequestJson(relativeUrl *url.URL, data []byte) ([]byte, error) {
	return self.doPostRequest(relativeUrl, data, "application/json")
}

// doPostRequest provides a generalized POST call to the REST API. Returns the response as
// a []byte slice, or an error if something goes wrong.
func (self *SawtoothClientTransportRest) doPostRequest(relativeUrl *url.URL, data []byte, contentType string) ([]byte, error) {
	fullUrl := self.URL.ResolveReference(relativeUrl)

	request, err := http.NewRequest(http.MethodPost, fullUrl.String(), bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}
	request.Header.Set("Accept", "application/json")
	request.Header.Set("Content-Type", contentType)

	response, err := self.HttpClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != 200 && response.StatusCode != 202 {
		transportError := NewSawtoothClientTransportRestError(response)
		err = &errors.SawtoothClientTransportError{
			ErrorCode:      errors.SawtoothTransportErrorCode(transportError.ErrorResponse.Error.Code),
			TransportError: transportError,
		}

		return nil, err
	}

	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	return responseData, nil
}
