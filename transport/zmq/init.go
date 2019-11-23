package zmq

import (
	"github.com/hyperledger/sawtooth-sdk-go/logging"
	"io/ioutil"
)

func init() {
	// Disable the logger from the sawtooth-sdk-go/logging package
	logger := logging.Get()
	logger.SetOutput(ioutil.Discard)
}
