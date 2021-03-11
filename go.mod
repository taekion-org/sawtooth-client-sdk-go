module github.com/taekion-org/sawtooth-client-sdk-go

go 1.16

replace github.com/hyperledger/sawtooth-sdk-go => github.com/taekion-org/sawtooth-sdk-go v0.1.4

require (
	github.com/brianolson/cbor_go v1.0.0
	github.com/golang/protobuf v1.4.3
	github.com/hyperledger/sawtooth-sdk-go v0.0.0-00010101000000-000000000000
	github.com/pebbe/zmq4 v1.2.5
	github.com/sirupsen/logrus v1.8.1
	github.com/spf13/pflag v1.0.5
)
