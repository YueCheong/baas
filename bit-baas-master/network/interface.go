package network

import (
	"bit-bass/deploy"
	"github.com/docker/docker/api/types"
	"github.com/docker/go-connections/nat"
)

const (
	ORDERERIMAGE = "hyperledger/fabric-orderer:1.4"
	PEERIMAGE    = "hyperledger/fabric-peer:1.4"
	FABRICNET    = "fabric_net"
	TOOLSIMAGE   = "hyperledger/fabric-tools:1.4"
)

type OperateIf interface {
	// 创建操作，加载配置
	Create() error
	// 移除操作
	Remove() error
	// 启动操作
	Start() error
	// 停止操作
	Stop() error
	// 连接到docker网络
	ConnectNet(net deploy.NetIf) error
}

type NodeIf interface {
	// 复用操作接口，用于操作docker容器
	OperateIf
	// 返回节点ID
	NodeName() string
	// 返回节点所在域
	NodeOrg() OrgIf
	// 审查容器状态
	Inspect() (types.ContainerJSON, error)
	// 打印容器日志
	PrintLog() (int64, error)
	// 获取容器的唯一标识
	ContainerID() string
	// 获得容器的名称
	ContainerName() string
	// 获得容器向外暴露的服务端口
	ServePort() nat.Port
}

type OrgIf interface {
	// 复用操作接口，用于批量操作组织内的节点
	OperateIf
	// 返回组织的名称
	Name() string
	// 返回组织所在域
	Domain() string
	// 返回组织的成员管理ID
	MSPID() string
}
