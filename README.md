Sawtooth Client SDK for Go
======================================

This repo contains the Sawtooth Client SDK for Go (developed by Taekion).

What is it for?
---------------

Client libraries that interface with Sawtooth transaction processors tend to repeat a lot of boilerplate
code for transaction/batch construction and handling, as well as for performing queries against the data
already committed to the blockchain. This library attempts to generalize much of that code in such a way
that lets application-specific client libraries avoid handling many of these implementation details.

What does it do?
----------------
- Handles transaction/batch/batchlist construction and signing.
- Provides a complete abstraction of the Sawtooth Validator interface.
    - Implements this as a generalized "transport" abstraction.
    - Two implementations of transport provided:
        - REST API
        - Direct ZMQ (0MQ) to the validator 

Getting Started
---------------

A application-specific client library needs to implement the following interface:

    type SawtoothClientImpl interface {
        GetFamilyName() string
        GetFamilyVersion() string
    
        EncodePayload(interface{}) ([]byte, error)
        DecodePayload([]byte, interface{}) error
    
        EncodeData(interface{}) ([]byte, error)
        DecodeData([]byte, interface{}) error
    
        GetPayloadInputAddresses(payload interface{}) []string
        GetPayloadOutputAddresses(payload interface{}) []string
        GetPayloadDependencies(payload interface{}) []string
    }


Once this is done, the library code can then construct an instance of the general library like so:

    import (
        "github.com/taekion-org/sawtooth-client-sdk-go"
        "github.com/taekion-org/sawtooth-client-sdk-go/transport"
    )
            
    type AppSpecificClient struct {
        *sawtooth-client-sdk-go.SawtoothClient
    }
    
    // Assume that this type implements the SawtoothClientImpl interface
    type AppSpecificClientImpl struct {}
    
    func NewAppSpecificClient(url string, keyFile string) (*AppSpecificClient, error) {
        args := &sawtooth_client_sdk_go.SawtoothClientArgs{
            URL: url,
            KeyFile: keyFile,
            TransportType: transport.TRANSPORT_REST,  // This can also be transport.TRANSPORT_ZMQ
            Impl: &AppSpecificClientImpl{},
        }
    
        sawtoothClient, err := sawtooth_client_sdk_go.NewClient(args)
        if err != nil {
            return nil, err
        }
    
        appSpecificClient := &AppSpecificClientImpl{
            SawtoothClient: sawtoothClient,
        }
    
        return appSpecificClient, nil
    }

At this point, the basic structure of the client is in place. Application-specific logic and functionality
can be implemented using the functions that the general library provides for executing transactions and queries.

Example
-------
For a more complete example, see the `examples/intkey` example. This provides a more-or-less complete re-implementation
of the `intkey` client utility using this SDK.
