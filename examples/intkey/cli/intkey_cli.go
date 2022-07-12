// intkey_cli is an alternative cli implementation for intkey that uses the example client implemented
// with sawtooth-client-sdk-go.
package main

import (
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"time"

	flag "github.com/spf13/pflag"
	"github.com/taekion-org/sawtooth-client-sdk-go/examples/intkey"
	"github.com/taekion-org/sawtooth-client-sdk-go/transport/errors"
)

const DEFAULT_REST_URL = "http://localhost:8008"
const DEFAULT_ZMQ_URL = "tcp://localhost:4004"
const DEFAULT_TRANSPORT = "rest"
const DEFAULT_WAIT_TIME = 5

const CMD_LIST = "list"
const CMD_SHOW = "show"
const CMD_SET = "set"
const CMD_INC = "inc"
const CMD_DEC = "dec"
const CMD_STATUS = "status"

var wait *uint = flag.Uint("wait", DEFAULT_WAIT_TIME, "Time to wait for commit")
var intkeyClient *intkey.IntkeyClient

func init() {
	// initialize RNG
	rand.Seed(time.Now().UnixNano())
}

func main() {
	var err error
	var rest_url *string = flag.String("rest_url", DEFAULT_REST_URL, "Sawtooth REST API URL")
	var zmq_url *string = flag.String("zmq_url", DEFAULT_ZMQ_URL, "Sawtooth ZMQ URL")
	var keyFile *string = flag.String("keyfile", "", "Sawtooth Private Key File")
	var transport *string = flag.String("transport", DEFAULT_TRANSPORT, "Sawtooth Transport")
	flag.Parse()

	if *transport == "rest" {
		intkeyClient, err = intkey.NewIntkeyClient(*rest_url, *keyFile)
		if err != nil {
			handleError(err)
		}
	} else if *transport == "zmq" {
		intkeyClient, err = intkey.NewIntkeyClientZmq(*zmq_url, *keyFile)
		if err != nil {
			handleError(err)
		}
	} else {
		handleError(fmt.Errorf("Invalid transport"))
	}

	if flag.NArg() == 0 {
		fmt.Printf("Usage: %s list|show|set|inc|dec|status [params] {--url [URL]} {--wait [wait_time]} {--transport [rest|zmq]}\n", os.Args[0])
		os.Exit(0)
	}

	switch flag.Arg(0) {
	case CMD_LIST:
		cmdList()
	case CMD_SHOW:
		cmdShow()
	case CMD_SET:
		cmdSet()
	case CMD_INC:
		cmdInc()
	case CMD_DEC:
		cmdDec()
	case CMD_STATUS:
		cmdStatus()
	default:
		cmdInvalid()
	}
}

func cmdList() {
	list, err := intkeyClient.List()
	if err != nil {
		handleError(err)
	}

	for key, value := range list {
		fmt.Printf("%s: %d\n", key, value)
	}
}

func cmdShow() {
	key := flag.Arg(1)
	if key == "" {
		handleError(fmt.Errorf("Error: command 'show' requires a parameter"))
	}

	value, err := intkeyClient.Show(key)
	if err != nil && err.(*errors.SawtoothClientTransportError).ErrorCode == errors.STATE_NOT_FOUND {
		handleError(fmt.Errorf("Error: key %s not found", key))
	} else if err != nil {
		handleError(err)
	}

	fmt.Printf("%s: %d\n", key, value)
}

func cmdSet() {
	if flag.NArg() < 3 {
		handleError(fmt.Errorf("Error: command 'set' requires two parameters"))
	}
	key := flag.Arg(1)
	value, err := strconv.ParseUint(flag.Arg(2), 10, 32)
	if err != nil {
		handleError(err)
	}

	existingValue, err := intkeyClient.Show(key)
	if err == nil {
		handleError(fmt.Errorf("Error: key %s already exists with value %d", key, existingValue))
	}

	batchId, err := intkeyClient.Set(key, uint(value), *wait)
	printBatchInfo(batchId)
}

func cmdInc() {
	if flag.NArg() < 3 {
		handleError(fmt.Errorf("Error: command 'inc' requires two parameters"))
	}
	key := flag.Arg(1)
	value, err := strconv.ParseUint(flag.Arg(2), 10, 32)
	if err != nil {
		handleError(err)
	}

	_, err = intkeyClient.Show(key)
	if err != nil && err.(*errors.SawtoothClientTransportError).ErrorCode == errors.STATE_NOT_FOUND {
		handleError(fmt.Errorf("Error: key %s not found", key))
	}

	batchId, err := intkeyClient.Inc(key, uint(value), *wait)
	printBatchInfo(batchId)

}

func cmdDec() {
	if flag.NArg() < 3 {
		handleError(fmt.Errorf("Error: command 'dec' requires two parameters"))
	}
	key := flag.Arg(1)
	value, err := strconv.ParseUint(flag.Arg(2), 10, 32)
	if err != nil {
		handleError(err)
	}

	_, err = intkeyClient.Show(key)
	if err != nil && err.(*errors.SawtoothClientTransportError).ErrorCode == errors.STATE_NOT_FOUND {
		handleError(fmt.Errorf("Error: key %s not found", key))
	}

	batchId, err := intkeyClient.Dec(key, uint(value), *wait)
	printBatchInfo(batchId)
}

func cmdStatus() {
	if flag.NArg() < 2 {
		handleError(fmt.Errorf("Error: command 'status' requires a parameter"))
	}
	batchId := flag.Arg(1)
	status, err := intkeyClient.Status(batchId)
	if err != nil {
		handleError(err)
	}

	fmt.Println(status)
}

func cmdInvalid() {
	handleError(fmt.Errorf("Error: '%s is an invalid command", flag.Arg(1)))
}

func printBatchInfo(batchId string) {
	status, err := intkeyClient.Status(batchId)
	if err != nil {
		handleError(err)
	}
	fmt.Printf("Batch: %s [Status: %s]\n", batchId, status)
}

func handleError(err error) {
	fmt.Println(err)
	os.Exit(-1)
}
