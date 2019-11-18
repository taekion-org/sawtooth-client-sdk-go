package zmq

import (
	"github.com/golang/protobuf/proto"
	"github.com/hyperledger/sawtooth-sdk-go/protobuf/validator_pb2"
)

func (self *SawtoothClientTransportZmq) doZmqRequest(t validator_pb2.Message_MessageType, request proto.Message, response proto.Message) error {
	requestMsg, err := proto.Marshal(request)
	if err != nil {
		return err
	}
	corrId, err := self.Connection.SendNewMsg(t, requestMsg)
	if err != nil {
		return err
	}

	_, responseMsg, err := self.Connection.RecvMsgWithId(corrId)
	err = proto.Unmarshal(responseMsg.GetContent(), response)
	if err != nil {
		return err
	}

	// TODO: We need code here to match common errors and generate error objects

	return nil
}
