package zutils

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/happylay-cloud/gf-extend/net/ztcp/ziface"
	"github.com/happylay-cloud/gf-extend/net/ztcp/zlog"
)

// 存储一切有关框架的全局参数，供其他模块使用，用户也可以通过根据tcp.json来配置
type GlobalObj struct {

	// 服务器配置
	TcpServer ziface.IServer // 当前gf-plus的全局Server对象
	Host      string         // 当前服务器主机IP
	TcpPort   int            // 当前服务器主机监听端口号
	Name      string         // 当前服务器名称

	// 高级配置
	Version          string // 版本号
	MaxPacketSize    uint32 // 都需数据包的最大值
	MaxConn          int    // 当前服务器主机允许的最大链接个数
	WorkerPoolSize   uint32 // 业务工作Worker池的数量
	MaxWorkerTaskLen uint32 // 业务工作Worker对应负责的任务队列最大任务存储数量
	MaxMsgChanLen    uint32 // SendBuffMsg发送消息的缓冲最大长度

	// 文件路径配置
	ConfFilePath string

	// 日志配置
	LogDir        string // 日志所在文件夹 默认"./log"
	LogFile       string // 日志文件名称   默认""，如果没有设置日志文件，打印信息将打印至stderr
	LogDebugClose bool   // 是否关闭Debug日志级别调试信息 默认false，默认打开debug信息
}

// 定义一个全局的对象
var GlobalObject *GlobalObj

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
func (global *GlobalObj) Reload() {

	if confFileExists, _ := PathExists(global.ConfFilePath); confFileExists != true {
		//fmt.Println("配置文件", global.ConfFilePath , "不存在！")
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

	// 日志设置
	if global.LogFile != "" {
		zlog.SetLogFile(global.LogDir, global.LogFile)
	}
	if global.LogDebugClose == true {
		zlog.CloseDebug()
	}
}

// 提供init方法，默认加载
func init() {
	pwd, err := os.Getwd()
	if err != nil {
		pwd = "."
	}
	// 初始化GlobalObject变量，设置一些默认值
	GlobalObject = &GlobalObj{
		Name:             "Gf-Plus TcpServerApp",
		Version:          "V1.1.0",
		TcpPort:          2021,
		Host:             "0.0.0.0",
		MaxConn:          12000,
		MaxPacketSize:    4096,
		ConfFilePath:     pwd + "/config/tcp.json",
		WorkerPoolSize:   10,
		MaxWorkerTaskLen: 1024,
		MaxMsgChanLen:    1024,
		LogDir:           pwd + "/log",
		LogFile:          "",
		LogDebugClose:    false,
	}

	// 从配置文件中加载一些用户配置的参数
	GlobalObject.Reload()
}
