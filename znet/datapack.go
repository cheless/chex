package znet

import (
	"bytes"
	"encoding/binary"
	"errors"
	"github.com/cheless/chex/utils"
	"github.com/cheless/chex/ziface"
)

type DataPack struct{}

// 实现 IDataPack 接口

func NewDataPack() *DataPack {
	return &DataPack{}
}

func (dp *DataPack) GetHeadLen() uint32 {
	return utils.Global.HeadLen
}

func (dp *DataPack) Pack(msg ziface.IMessage) ([]byte, error) {
	// 创建一个存放 []byte 的缓冲
	dataBuff := bytes.NewBuffer([]byte{})

	// 将 DataLen 写入
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetDataLen()); err != nil {
		return nil, err
	}
	// 将 MsgID 写入
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetMsgID()); err != nil {
		return nil, err
	}
	// 将 Data 写入
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetData()); err != nil {
		return nil, err
	}
	return dataBuff.Bytes(), nil
}

func (dp *DataPack) Unpack(data []byte) (ziface.IMessage, error) {
	// 创建读取缓冲
	dataBuff := bytes.NewReader(data)

	// 将 Message head 的各部分分别读入，因为 head 长度固定，DataLen 和 ID 都是 uint32 类型
	msg := &Message{}
	if err := binary.Read(dataBuff, binary.LittleEndian, &msg.DataLen); err != nil {
		return nil, err
	}
	if err := binary.Read(dataBuff, binary.LittleEndian, &msg.Id); err != nil {
		return nil, err
	}

	// 判断 DateLen 长度是否大于自定义允许的最大包长度
	if utils.Global.MaxPacketSize > 0 && utils.Global.MaxPacketSize < msg.DataLen {
		return nil, errors.New("Message data receive is too beyound MaxPacketSize")
	}

	return msg, nil
}
