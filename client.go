// Package sawtooth_client_sdk_go provides an easy-to-use SDK for building Sawtooth client libraries.
package sawtooth_client_sdk_go

import (
	"fmt"
	"github.com/hyperledger/sawtooth-sdk-go/signing"
	"github.com/taekion-org/sawtooth-client-sdk-go/transport"
	"net/url"
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
	PrivateKey		signing.PrivateKey
	KeyFile			string
	Impl			SawtoothClientImpl
	TransportType	transport.SawtoothClientTransportType
}

// NewClient constructs a new instance of the SawtoothClient.
func NewClient(args *SawtoothClientArgs) (*SawtoothClient, error) {
	var err error
	var privateKey signing.PrivateKey

	// Set up our private key for signing.
	// If we are passed a private key directly, use it.
	// Otherwise, load one from a file.
	if args.PrivateKey != nil {
		privateKey = args.PrivateKey
	} else {
		var keyFile string
		if args.KeyFile == "" {
			keyFile, err = getDefaultKeyFileName()
			if err != nil {
				return nil, err
			}
		} else {
			keyFile = args.KeyFile
		}

		privateKey, err = readPrivateKeyFromFile(keyFile)
		if err != nil {
			return nil, err
		}
	}

	// Set up the crypto factory and signer
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
