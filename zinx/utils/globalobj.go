package utils

import (
	"encoding/json"
	"io/ioutil"
	"zinx-test/zinx/ziface"
)

/*
存储一切有关Zinx框架的全局参数，供其他模块使用
一些参数是可以通过zinx》json由用户进行配置
*/

type GlobalObj struct {
	/*
		Server
	*/
	TcpServer ziface.IServer
	Host      string
	TcpPort   int
	Name      string
	/*
		Zinx
	*/
	Version        string // 当前zinx的version
	MaxConn        int    // 当前服务器主机允许的最大链接数
	MaxPackageSize uint32
}

/*
定义一个全局的对外Globalobj
*/
var GlobalObject *GlobalObj

/*
从 zinx.json去加载用于自定义的参数
*/

func (g *GlobalObj) Reload() {
	data, err := ioutil.ReadFile("conf/zinx.json")
	if err != nil {
		panic(err)
	}
	// 将json文件数据解析到struct中
	err = json.Unmarshal(data, &GlobalObject)
	if err != nil {
		panic(err)
	}
}

/*
提供一个init方法，初始化当前的GlobalObj
*/
func init() {
	// 如果配置文件没有加载，默认的值
	GlobalObject = &GlobalObj{
		Name:           "ZinxServerApp",
		Version:        "V0.7",
		TcpPort:        8999,
		Host:           "0.0.0.0",
		MaxConn:        1000,
		MaxPackageSize: 4096,
	}
	// 应该尝试从conf/zinx.json去加载一些用户配置的入参
	GlobalObject.Reload()
}
