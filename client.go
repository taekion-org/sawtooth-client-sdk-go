// Package sawtooth_client_sdk_go provides an easy-to-use SDK for building Sawtooth client libraries.
package sawtooth_client_sdk_go

import (
	"fmt"
	"github.com/hyperledger/sawtooth-sdk-go/signing"
	"github.com/taekion-org/sawtooth-client-sdk-go/transport"
	"io/ioutil"
	"net/url"
	"strings"
)

// SawtoothClient represents the core functionality of a Sawtooth application client.
type SawtoothClient struct {
	Signer			*signing.Signer
	Transport		transport.SawtoothClientTransport
	ClientImpl		SawtoothClientImpl
}

// SawtoothClientArgs holds arguments required to initialize SawtoothClient.
type SawtoothClientArgs struct {
	URL				string
	KeyFile			string
	Impl			SawtoothClientImpl
	TransportType	transport.SawtoothClientTransportType
}

// NewClient constructs a new instance of the SawtoothClient.
func NewClient(args *SawtoothClientArgs) (*SawtoothClient, error) {
	var err error

	// Figure out which key to use
	var keyFile string
	if args.KeyFile == "" {
		keyFile, err = getDefaultKeyFileName()
		if err != nil {
			return nil, err
		}
	} else {
		keyFile = args.KeyFile
	}

	// Read the key
	keyData, err := ioutil.ReadFile(keyFile)
	if err != nil {
		return nil, fmt.Errorf("Could not read private key from file (%s) with error: %s", keyFile, err)
	}
	keyData = []byte(strings.TrimSpace(string(keyData)))

	// Set up the key structures and signer
	privateKey := signing.NewSecp256k1PrivateKey(keyData)
	cryptoFactory := signing.NewCryptoFactory(signing.CreateContext("secp256k1"))
	signer := cryptoFactory.NewSigner(privateKey)

	// Parse the URL
	url, err := url.Parse(args.URL)
	if err != nil {
		return nil, fmt.Errorf("Error parsing URL: %s", err)
	}

	// Create the transport
	transport, err := transport.NewSawtoothClientTransport(args.TransportType, url)
	if err != nil {
		return nil, fmt.Errorf("Error initializing transport: %s", err)
	}

	client := &SawtoothClient{Signer: signer, ClientImpl: args.Impl, Transport: transport}

	return client, nil
}
