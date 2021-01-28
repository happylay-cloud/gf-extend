package hutils

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/happylay-cloud/gf-extend/net/htcp/hiface"
)

// HTcpGlobalConfig 全局配置
type HTcpGlobalConfig struct {
	TcpServer        hiface.ITcpServer // 当前gf-plus的全局TcpServer对象
	Host             string            // 当前服务器主机IP
	TcpPort          int               // 当前服务器主机监听端口号
	Name             string            // 当前服务器名称
	Version          string            // 版本号
	MaxPacketSize    int64             // 接收数据包的最大值
	MaxPkgChanLen    int64             // 发送消息的缓冲最大长度
	MaxConn          int               // 当前服务器主机允许的最大连接个数
	WorkerPoolSize   int64             // 业务工作Worker池的数量
	MaxWorkerTaskLen int64             // 业务工作Worker对应负责的任务队列最大任务存储数量
	ConfFilePath     string            // 配置文件路径
}

// 定义一个全局的对象
var GlobalHTcpObject *HTcpGlobalConfig

// 判断一个文件是否存在
func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

// 读取用户的配置文件
func (global *HTcpGlobalConfig) Reload() {

	if confFileExists, _ := PathExists(global.ConfFilePath); confFileExists != true {
		return
	}

	data, err := ioutil.ReadFile(global.ConfFilePath)
	if err != nil {
		panic(err)
	}
	// 将json数据解析到struct中
	err = json.Unmarshal(data, global)
	if err != nil {
		panic(err)
	}

}

// 初始化全局变量
func init() {
	pwd, err := os.Getwd()
	if err != nil {
		pwd = "."
	}
	// 初始化GlobalHTcpObject变量，设置一些默认值
	GlobalHTcpObject = &HTcpGlobalConfig{
		Name:             "Gf-Plus TcpServerApp",
		Version:          "V1.1.0",
		Host:             "0.0.0.0",
		TcpPort:          2021,
		MaxConn:          12000,
		WorkerPoolSize:   10,
		MaxPacketSize:    4096,
		MaxWorkerTaskLen: 1024,
		MaxPkgChanLen:    1024,
		ConfFilePath:     pwd + "/config/tcp.json",
	}
	// 从配置文件中加载一些用户配置的参数
	GlobalHTcpObject.Reload()
}
