package rest

import (
	"bytes"
	"fmt"
	"github.com/taekion-org/sawtooth-client-sdk-go/transport/errors"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"
)

// resolveReference takes the base URL and combines it with the specified relativeUrl.
// This works differently from the standard URL.ResolveReference() in that it always
// adds the relativeUrl to the end of the base URL in the "right way".
func (self *SawtoothClientTransportRest) resolveReference(relativeUrl *url.URL) *url.URL {
	newUrl := *self.URL
	newUrl.Path = path.Join(newUrl.Path, relativeUrl.Path)
	newUrl.RawQuery = relativeUrl.RawQuery

	return &newUrl
}

// buildRequest wraps http.NewRequest and sets up the headers as we require them.
// In particular we set Accept to "application/json", and check the URL for authentication
// information. If the username specified is "bearer", we add an Authorization header set to
// "Bearer <value_of_password_field>".
func (self *SawtoothClientTransportRest) buildRequest(method string, url *url.URL, body io.Reader) (*http.Request, error){
	urlString := url.String()

	request, err := http.NewRequest(method, urlString, body)
	if err != nil {
		return nil, err
	}
	request.Header.Set("Accept", "application/json")

	authMethod := url.User.Username()
	authSecret, authSecretPresent := url.User.Password()

	if authMethod == "bearer" && authSecretPresent {
		request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", authSecret))
	}

	return request, nil
}

// doGetRequest provides a generalized GET call to the REST API. Returns the response as
// a []byte slice, or an error if something goes wrong.
func (self *SawtoothClientTransportRest) doGetRequest(relativeUrl *url.URL) ([]byte, error) {
	fullUrl := self.resolveReference(relativeUrl)

	request, err := self.buildRequest(http.MethodGet, fullUrl, nil)
	if err != nil {
		return nil, err
	}

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
	fullUrl := self.resolveReference(relativeUrl)

	request, err := self.buildRequest(http.MethodPost, fullUrl, bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}
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
