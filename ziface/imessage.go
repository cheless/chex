package ziface

// 将请求的消息封装到 Message 结构中

type IMessage interface {
	GetMsgID() uint32
	GetDataLen() uint32
	GetData() []byte

	SetMsgID(uint32)
	SetMsgLen(uint32)
	SetData([]byte)
}
