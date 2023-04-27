package deploy

import (
	"github.com/docker/docker/api/types"
	"github.com/docker/go-connections/nat"
)

type NodeConfigIf interface {
	// 节点的主机名
	Host() string
	// 节点所在域
	Domain() OrgConfigIf
	// 节点的服务端口
	Port() nat.Port
	// 节点加入组织
	JoinOrg(domain OrgConfigIf) error
	// 节点链接的网络名
	NetName() string
}

type OrgConfigIf interface {
	// 返回组织的名称
	Name() string
	// 返回组织所在域
	Domain() string
	// 返回组织的成员管理ID
	MSPID() string
	// 返回成员管理密钥目录
	CryptoPath() string
	// 返回创世区块、通道交易所在路径
	ConfigtxPath() string
	// 向组织中添加节点
	AddNode(node NodeConfigIf) error
	// 根据host移除该节点
	DelNode(host string) error
	// 获得当前组织中的所有节点
	GetNodes() []NodeConfigIf
}

type NetIf interface {
	// 返回网络名称
	NetName() string
	// 返回网络ID
	NetID() string
	// 创建docker网络
	CreateNet() error
	// 移除docker网络
	RemoveNet() error
	// 判断当前网络是否存在
	IsNetExist() bool
	// 审查docker网络状态
	Inspect() types.NetworkResource
}
