package utils

import (
	"encoding/json"
	"io/ioutil"
	"github.com/cheless/chex/ziface"
)

/*
	存储一切有关Zinx框架的全局参数，供其他模块使用
	用户可以根据 zinx.json 来配置其中一些参数
*/
type GlobalObj struct {
	TcpServer ziface.IServer // 当前Zinx的全局Server对象
	Host      string         // 当前服务器主机IP
	TcpPort   int            // 当前服务器主机监听端口号
	Name      string         // 当前服务器名称
	Version   string         // 当前Zinx版本号

	MaxPacketSize  uint32 // 读取数据包的最大值
	MaxConn        int    // 当前服务器主机允许的最大链接个数
	HeadLen        uint32 // 封包的后包 head 的长度
	WorkerPoolSize uint32 // WorkerPool 中 worker 的数量
	MaxWorkerTask  uint32 // 每个 worker 对应的消息队列的最大任务数量
}

var Global *GlobalObj

// 读取用户的配置文件
func (g *GlobalObj) Reload() {
	data, err := ioutil.ReadFile("../conf/zinx.json")
	if err != nil {
		panic(err)
	}
	// 将json数据解析到struct中
	// fmt.Printf("json :%s\n", data)
	err = json.Unmarshal(data, &Global)
	if err != nil {
		panic(err)
	}
}

/*
	提供init方法，默认加载
*/
func init() {
	// 初始化默认值
	Global = &GlobalObj{
		Name:           "ZinxServerApp",
		Version:        "V0.9",
		TcpPort:        7777,
		Host:           "0.0.0.0",
		MaxConn:        12000,
		MaxPacketSize:  4096,
		HeadLen:        8, // DateLen(uint32) + ID(uint32) = 8B
		WorkerPoolSize: 8,
		MaxWorkerTask:  1024,
	}

	// 将用户传入的参数更新到 GlobalObj
	Global.Reload()
}
