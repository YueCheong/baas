package core

import (
	"bit-bass/contract"
)

//用于获取网络和创建网络
type BlockchainToJson struct {
	//每个区块链网络都有一个独一无二的ID
	//在创建网络的请求中此字段无用
	ID int
	//区块链网络的名称，同一个用户的网络名称不能重复
	Name string
	//区块链网络中的Orderer组织，目前只有一个.
	OrdererOrg OrdererOrgToJson
	//区块链网络中的Peer组织，可以有多个
	PeerOrg []PeerOrgToJson
	//在创建区块链网络过程中可选的额外参数，留作以后扩展。可以留空，代表采用默认设置
	Config map[string]interface{}
	//区块链网络中的通道
	Channels []ChannelToJson
	//区块连系统的状态
	Status string
	//区块链系统链接的网络名称
	Netname string
}

type BlockchainOperationType int

const (
	StartBlockchain BlockchainOperationType = iota + 1
	StopBlockchain
	InitializeBlockchain
	SetOrdererOrg
	SetOrderer
	AddPeerOrg
	AddPeer
)

type BlockchainManageInfoToJson struct {
	BlockchainID int

	Operation BlockchainOperationType

	Name string

	Domain string

	MSPID string

	Nodes []NodeToJson
}

type OrdererOrgToJson struct {
	//注意，Orderer组织的名字加域名不能重复
	//Orderer组织的名字
	Name string
	//Orderer组织的域名
	Domain string
	//Orderer组织中的Orderer
	Orderer NodeToJson
	//Orderer组织的MSPID
	MSPID string
}

type PeerOrgToJson struct {
	//注意，Peer组织的名字加域名不能重复
	//Peer组织的名字，如"org1"
	Name string
	//Peer组织的域名,
	//例子："example.com"
	Domain string
	//Peer组织的MSPID
	MSPID string
	//Peer组织中的节点们
	Peers []NodeToJson
}

type NodeToJson struct {
	//Peer节点的名字，同一个组织中的节点不能重名
	Name string
	//Peer节点使用的接口，同一个网络中的节点不能重复
	Port string
}

//用于获取和创建通道
type ChannelToJson struct {
	//通道的名字，同一个区块链网络中不能有重复的通道名。即在一个区块链网络中
	//可以用通道名唯一的确定一个通道
	Name string
	//通道中的Peer节点，通过Peer节点名加Peer节点所在的组织名以及域名确定
	//例子： "peer0.org1.example.com"
	Peers []string
	//该通道中为每个peer组织设定的锚节点，键名为组织名加域名，值为peer
	//节点的全名，即peer名加组织名加域名
	//例子：  "org1.example.com" : "peer0.org1.example.com"
	AnchorPeers map[string]string
	//通道所属区块链网络的ID
	BlockchainID int
	//通道所属区块链网络的名字
	Blockchainname string

	//创建通道的组织
}

//用于改变通道的状态，如添加节点，改变锚节点等
type ChannelManageInfoToJson struct {
	//通道所属的区块链网络ID
	BlockchainID int
	//通道所属的区块链网络的名字
	BlockchainName string
	//通道的名字，用于制定操作的对象通道
	ChannelName string
	//指明对通道进行的操作，类型为int
	Operation ChannelOperationType
	//对通道进行操作时所需的参数
	Args []string
}

type ChannelOperationType int

const (
	//向通道中添加节点
	AddPeerToChannel ChannelOperationType = iota
	//为通道中的某个组织设定锚节点
	SetAnchorPeer
)

func (t ChannelOperationType) String() string {
	switch t {
	case AddPeerToChannel:
		return "AddPeerToChannel"
	case SetAnchorPeer:
		return "SetAnchorPeer"
	default:
		return "Unkown Operation"
	}
}

//用于获取和创建合约
type ContractToJson struct {
	//合约的ID，每个合约都拥有一个独一无二的ID，由系统自动分配
	//在创建合约的请求中此字段无用
	ID int
	//合约所属的区块链网络ID
	BlockchainID int
	//合约所属的区块链网络名字
	BlockchainName string
	//合约的名字
	ContractName string
	//合约所在的通道名
	ChannelName string
	//合约的描述
	ContractDesc string
	//合约的版本
	ContractVersion string
	//合约的路径
	ContractPath string
	//合约使用的编程语言
	ContractLang contract.ChaincodeLanguageType
}

//用于合约状态的改变，实例化及升级合约
type ContractManageInfoToJson struct {
	//需要操作的合约ID
	ID int
	//需要操作的合约所在区块链的ID
	BlockchainID int
	//需要操作的合约名
	ContractName string
	//需要操作的合约版本
	ContractVersion string
	//需要对合约进行的操作
	Operation ContractOperationType
	//对合约操作的参数
	Args string
}

type ContractOperationType int

const (
	Install ContractOperationType = iota
	Instantiate
	Upgrade
)

func (t ContractOperationType) String() string {
	switch t {
	case Install:
		return "Install"
	case Instantiate:
		return "Instantiate"
	case Upgrade:
		return "Upgrade"
	default:
		return "Unknown"
	}
}

//用于查询和调用智能合约
type ContractInvokeInfoToJson struct {
	//需要调用的合约ID
	ID int
	//需要调用的合约所在的区块链网络ID
	BlockchainID int
	//需要调用的合约名
	ContractName string
	//需要调用的合约版本
	ContractVersion string
	//需要进行的调用操作类型，Query或Invoke
	InvokeType ContractInvokeType
	//调用操作的参数
	Args string
}

//合约调用的结果
type CCInvokeResponseToJson struct {
	//交易id
	Txid string
	//合约调用的状态码
	Status int32
	//合约调用的返回结果
	Payload string
}

type ContractInvokeType int

const (
	Query = iota
	Invoke
)

type BaasSummaryToJson struct {
	BlochainTotal      int
	BlockchainRunning  int
	OrganizationTotal  int
	OrdererOrgTotal    int
	PeerOrgTotal       int
	ChannelTotal       int
	ContractTotal      int
	BlockTotal         string
	BlockchainSummarys []BlockchainSummary
}

type BlockchainSummary struct {
	BlockchainName  string
	BlockchainID    int
	OrdererOrgNum   int
	PeerOrgNum      int
	ChannelNum      int
	ContractNum     int
	ChannelSummarys []ChannelSummary
}

type ChannelSummary struct {
	ChannelName string
	Blockheight string
}
