package contract

import (
	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/resmgmt"
)

type ContractIF interface {
	// 返回链码ID
	ChaincodeID() string
	// 返回链码版本号
	ChaincodeVer() string
	// 返回链码描述
	ChaincodeDesc() string
	// 返回链码所在通道
	ChaincodeChan() string
	// 安装链码
	InstallChaincode() ([]resmgmt.InstallCCResponse, error)
	// 初始化链码
	InitialChaincode() (resmgmt.InstantiateCCResponse, error)
	// 调用链码
	InvokeChaincode(fcname string, args ...string) (channel.Response, error)
	// 查询链码
	QueryChaincode(fcname string, args ...string) (channel.Response, error)
}
