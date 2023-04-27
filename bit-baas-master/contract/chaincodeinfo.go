package contract

import (
	"bit-bass/network"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
	"github.com/pkg/errors"
)

//创建链码对象所需的参数
type ChaincodeConfig struct {
	ChannelName      string
	ChaincodeName    string
	ChaincodeDesc    string
	ChaincodeGoPath  string
	ChaincodePath    string
	ChaincodeVersion string
	ChaincodeLang    ChaincodeLanguageType
}

//链码对象
type ChaincodeInfo struct {
	id               int
	configPath       string
	channelName      string
	chaincodeName    string
	chaincodeDesc    string
	chaincodeGoPath  string
	chaincodePath    string
	chaincodeVersion string
	chaincodeLang    ChaincodeLanguageType
}

func (ci *ChaincodeInfo) ChaincodeID() int {
	return ci.id
}

func (ci *ChaincodeInfo) ChaincodeDesc() string {
	return ci.chaincodeDesc
}

func (ci *ChaincodeInfo) ChaincodeVer() string {
	return ci.chaincodeVersion
}

func (ci *ChaincodeInfo) ChaincodeLanguage() ChaincodeLanguageType {
	return ci.chaincodeLang
}

func (ci *ChaincodeInfo) ChaincodeName() string {
	return ci.chaincodeName
}

func (ci *ChaincodeInfo) ChannelName() string {
	return ci.channelName
}

//创建链码管理客户端
func (ci *ChaincodeInfo) NewChaincodeClient(user string, org *network.PeerOrg) (*ChaincodeClient, error) {
	netConifg := config.FromFile(ci.configPath + "/sdk-config.yaml")
	fabricsdk, err := fabsdk.New(netConifg)
	if err != nil {
		return nil, errors.Errorf("Faild to creat fabric sdk : %v", err)
	}

	ccCli := ChaincodeClient{
		info:        ci,
		sdk:         fabricsdk,
		operateUser: user,
		operateOrg:  org,
	}

	return &ccCli, nil
}
