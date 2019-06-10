package znet

import "github.com/TonyXMH/ZinxDemo/ziface"

type Message struct {
	DataLen uint32
	ID      uint32
	Data    []byte
}

func (m *Message) GetDataLen() uint32 {
	return m.DataLen
}
func (m *Message) GetMsgID() uint32 {
	return m.ID
}
func (m *Message) GetData() []byte {
	return m.Data
}

func (m *Message) SetDataLen(l uint32) {
	m.DataLen = l
}
func (m *Message) SetMsgID(id uint32) {
	m.ID = id
}
func (m *Message) SetData(data []byte) {
	m.Data = data
}

func NewMsgPacket(id uint32, data []byte) ziface.IMessage {
	return &Message{
		DataLen: uint32(len(data)),
		ID:      id,
		Data:    data,
	}
}
