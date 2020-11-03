package zmq

import (
	"github.com/golang/protobuf/proto"
	"github.com/hyperledger/sawtooth-sdk-go/protobuf/validator_pb2"
	"github.com/taekion-org/sawtooth-client-sdk-go/transport/errors"
)

func (self *SawtoothClientTransportZmq) doZmqRequest(t validator_pb2.Message_MessageType, request proto.Message, response proto.Message) error {
	requestMsg, err := proto.Marshal(request)
	if err != nil {
		return errors.NewSawtoothClientTransportRequestError(err)
	}

	connection, err := self.getConnection()
	if err != nil {
		return errors.NewSawtoothClientTransportRequestError(err)
	}

	corrId, err := connection.SendNewMsg(t, requestMsg)
	if err != nil {
		return errors.NewSawtoothClientTransportRequestError(err)
	}

	_, responseMsg, err := connection.RecvMsgWithId(corrId)
	err = proto.Unmarshal(responseMsg.GetContent(), response)
	if err != nil {
		return errors.NewSawtoothClientTransportRequestError(err)
	}

	err = self.putConnection(connection)
	if err != nil {
		return errors.NewSawtoothClientTransportRequestError(err)
	}

	errorCode := checkForError(response)
	if errorCode != errors.NO_ERROR {
		transportError := NewSawtoothClientTransportZmqError(t, request, response, errorCode)
		err = &errors.SawtoothClientTransportError{
			ErrorCode: transportError.ErrorCode,
			ErrorObject: transportError,
		}
		return err
	}

	return nil
}
