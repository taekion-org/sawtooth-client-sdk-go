package rest

import (
	"github.com/golang/protobuf/proto"
	"github.com/hyperledger/sawtooth-sdk-go/protobuf/batch_pb2"
	"net/url"
)

// SubmitBatchList submits a batch list to Sawtooth. The batch list must be in the form of a
// batch_pb2.BatchList protobuf and be prepared appropriately (all required fields and signatures
// populated.
func (self *SawtoothClientTransportRest) SubmitBatchList(batchList *batch_pb2.BatchList) error {
	batchesSerialized, err := proto.Marshal(batchList)
	if err != nil {
		return err
	}

	relativeUrl := &url.URL{Path: "/batches"}
	_, err = self.doPostRequestBinary(relativeUrl, batchesSerialized)
	if err != nil {
		return err
	}

	return nil
}
