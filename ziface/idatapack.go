package ziface

// 封包和拆包模块，用于处理 TCP 黏包问题

type IDataPack interface {
	GetHeadLen() uint32
	Pack(msg IMessage) ([]byte, error)
	Unpack(data []byte) (IMessage, error)
}
