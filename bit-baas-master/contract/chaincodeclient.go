package contract

import (
	"bit-bass/network"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/resmgmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/errors/retry"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/fab/ccpackager/gopackager"
	"github.com/hyperledger/fabric-sdk-go/pkg/fab/ccpackager/javapackager"
	"github.com/hyperledger/fabric-sdk-go/pkg/fab/ccpackager/nodepackager"
	"github.com/hyperledger/fabric-sdk-go/pkg/fab/resource"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
	"github.com/hyperledger/fabric-sdk-go/third_party/github.com/hyperledger/fabric/common/cauthdsl"
	"github.com/pkg/errors"
)

type ChaincodeClient struct {
	info        *ChaincodeInfo
	sdk         *fabsdk.FabricSDK
	operateUser string
	operateOrg  *network.PeerOrg
}

// 返回链码ID
func (c *ChaincodeClient) ChaincodeID() int {
	return c.info.ChaincodeID()
}

// 返回链码名
func (c *ChaincodeClient) ChaincodeName() string {
	return c.info.ChaincodeName()
}

// 返回链码版本号
func (c *ChaincodeClient) ChaincodeVer() string {
	return c.info.ChaincodeVer()
}

// 返回链码描述
func (c *ChaincodeClient) ChaincodeDesc() string {
	return c.info.ChaincodeDesc()
}

// 返回链码所在通道
func (c *ChaincodeClient) ChaincodeChan() string {
	return c.info.ChannelName()
}

// 返回链码所使用的语言
func (c *ChaincodeClient) ChaincodeLanaguage() ChaincodeLanguageType {
	return c.info.ChaincodeLanguage()
}

func (c *ChaincodeClient) InstallChaincode() ([]resmgmt.InstallCCResponse, error) {
	var ccPkg *resource.CCPackage
	var err error = nil
	switch c.info.chaincodeLang {
	case Golang:
		ccPkg, err = gopackager.NewCCPackage(c.info.chaincodePath, c.info.chaincodeGoPath)
	case Java:
		ccPkg, err = javapackager.NewCCPackage(c.info.chaincodePath)
	case Node:
		ccPkg, err = nodepackager.NewCCPackage(c.info.chaincodePath)
	default:
		return nil, errors.New("Faild to pack chaincode : Unknown chaincode language ")
	}
	if err != nil {
		return nil, errors.Errorf("Failed to pack chaincode : %v", err)
	}

	installCCReq := resmgmt.InstallCCRequest{
		Name:    c.info.chaincodeName,
		Path:    c.info.chaincodePath,
		Version: c.info.chaincodeVersion,
		Package: ccPkg,
	}

	resClient, err := resmgmt.New(c.sdk.Context(fabsdk.WithUser(c.operateUser), fabsdk.WithOrg(c.operateOrg.Name())))
	if err != nil {
		return nil, errors.Errorf("Failed to creat resmgmt client : %v", err)
	}

	resp, err := resClient.InstallCC(installCCReq, resmgmt.WithRetry(retry.DefaultResMgmtOpts))

	if err != nil {
		return nil, errors.Errorf("Failed to install CC : %v", err)
	}

	return resp, nil
}

func (c *ChaincodeClient) InstantiateChaincode(args ...string) (resmgmt.InstantiateCCResponse, error) {
	var packedArgs [][]byte
	for _, arg := range args {
		packedArgs = append(packedArgs, []byte(arg))
	}

	ccPolicy := cauthdsl.SignedByMspMember(c.operateOrg.MSPID())
	initialCCReq := resmgmt.InstantiateCCRequest{
		Name:    c.info.chaincodeName,
		Path:    c.info.chaincodePath,
		Version: c.info.chaincodeVersion,
		Args:    packedArgs,
		Policy:  ccPolicy,
	}

	resClient, err := resmgmt.New(c.sdk.Context(fabsdk.WithUser(c.operateUser), fabsdk.WithOrg(c.operateOrg.Name())))
	if err != nil {
		return resmgmt.InstantiateCCResponse{}, errors.Errorf("Failed to creat resmgmt client : %v", err)
	}

	resp, err := resClient.InstantiateCC(c.info.channelName, initialCCReq, resmgmt.WithRetry(retry.DefaultResMgmtOpts))
	if err != nil {
		return resmgmt.InstantiateCCResponse{}, errors.Errorf("Failed to instantiate CC : %v", err)
	}
	return resp, nil
}

func (c *ChaincodeClient) UpgradeChaincode(args ...string) (resmgmt.UpgradeCCResponse, error) {
	var packedArgs [][]byte
	for _, arg := range args {
		packedArgs = append(packedArgs, []byte(arg))
	}

	ccPolicy := cauthdsl.SignedByMspMember(c.operateOrg.MSPID())
	upgradeCCReq := resmgmt.UpgradeCCRequest{
		Name:    c.info.chaincodeName,
		Path:    c.info.chaincodePath,
		Version: c.info.chaincodeVersion,
		Args:    packedArgs,
		Policy:  ccPolicy,
	}

	resClient, err := resmgmt.New(c.sdk.Context(fabsdk.WithUser(c.operateUser), fabsdk.WithOrg(c.operateOrg.Name())))
	if err != nil {
		return resmgmt.UpgradeCCResponse{}, errors.Errorf("Failed to creat resmgmt client : %v", err)
	}

	resp, err := resClient.UpgradeCC(c.info.channelName, upgradeCCReq, resmgmt.WithRetry(retry.DefaultResMgmtOpts))
	if err != nil {
		return resmgmt.UpgradeCCResponse{}, errors.Errorf("Failed to upgrade CC : %v", err)
	}
	return resp, nil
}

func (c *ChaincodeClient) InvokeChaincode(fcname string, args ...string) (channel.Response, error) {
	sdkChannelContext := c.sdk.ChannelContext(c.info.channelName,
		fabsdk.WithOrg(c.operateOrg.Name()), fabsdk.WithUser(c.operateUser))
	chClient, err := channel.New(sdkChannelContext)
	if err != nil {
		return channel.Response{}, errors.Errorf("Failed to creat sdk channel client : %v", err)
	}

	var packedArgs [][]byte
	for _, arg := range args {
		packedArgs = append(packedArgs, []byte(arg))
	}

	invokeRequest := channel.Request{
		ChaincodeID: c.info.chaincodeName,
		Fcn:         fcname,
		Args:        packedArgs,
	}

	var peernames []string
	for _, peer := range c.operateOrg.GetPeerNodes() {
		peernames = append(peernames, peer.ContainerName())
	}

	resp, err := chClient.Execute(invokeRequest)
	if err != nil {
		return channel.Response{}, errors.Errorf("Failed to execute invoke request : %v", err)
	}

	return resp, nil
}

func (c *ChaincodeClient) QueryChaincode(fcname string, args ...string) (channel.Response, error) {

	sdkChannelContext := c.sdk.ChannelContext(c.info.channelName,
		fabsdk.WithOrg(c.operateOrg.Name()), fabsdk.WithUser(c.operateUser))
	chClient, err := channel.New(sdkChannelContext)
	if err != nil {
		return channel.Response{}, errors.Errorf("Failed to creat sdk channel client : %v", err)
	}

	var packedArgs [][]byte
	for _, arg := range args {
		packedArgs = append(packedArgs, []byte(arg))
	}

	queryRequest := channel.Request{
		ChaincodeID: c.info.chaincodeName,
		Fcn:         fcname,
		Args:        packedArgs,
	}

	var peernames []string
	for _, peer := range c.operateOrg.GetPeerNodes() {
		peernames = append(peernames, peer.ContainerName())
	}

	resp, err := chClient.Query(queryRequest)
	if err != nil {
		return channel.Response{}, errors.Errorf("Failed to execute invoke request : %v", err)
	}

	return resp, nil
}

func (c *ChaincodeClient) Close() {
	c.sdk.Close()
}

//配置文件更新Channel对象中到sdk
func (c *ChaincodeClient) RefreshSdk() error {
	c.sdk.Close()
	netConifg := config.FromFile(c.info.configPath + "/sdk-config.yaml")
	fabricsdk, err := fabsdk.New(netConifg)
	if err != nil {
		return errors.Errorf("Faild to creat fabric sdk : %v", err)
	}

	c.sdk = fabricsdk
	return nil
}
