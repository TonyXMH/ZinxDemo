package znet

import (
	"bytes"
	"encoding/binary"
	"github.com/TonyXMH/ZinxDemo/utils"
	"github.com/TonyXMH/ZinxDemo/ziface"
	"github.com/pkg/errors"
)

type DataPack struct{}

func NewDataPack() ziface.IDataPack {
	return &DataPack{}
}

func (d *DataPack) GetHeadLen() uint32 {
	return 8 //uint32+uint32 4+4 = 8
}

func (d *DataPack) Pack(msg ziface.IMessage) (data []byte, err error) {
	buff := bytes.NewBuffer([]byte{})
	if err := binary.Write(buff, binary.LittleEndian, msg.GetDataLen()); err != nil {
		return nil, err
	}

	if err := binary.Write(buff, binary.LittleEndian, msg.GetMsgID()); err != nil {
		return nil, err
	}
	if err := binary.Write(buff, binary.LittleEndian, msg.GetData()); err != nil {
		return nil, err
	}

	return buff.Bytes(), nil
}

func (d *DataPack) Unpack(data []byte) (ziface.IMessage, error) {
	dataBuf := bytes.NewReader(data)
	msg := &Message{}
	if err := binary.Read(dataBuf, binary.LittleEndian, &msg.DataLen); err != nil {
		return nil, err
	}
	if err := binary.Read(dataBuf, binary.LittleEndian, &msg.ID); err != nil {
		return nil, err
	}

	if utils.GlobalObject.MaxPacketSize > 0 && msg.DataLen > utils.GlobalObject.MaxPacketSize {
		return nil, errors.New("Too Large message received")
	}
	return msg, nil
}
