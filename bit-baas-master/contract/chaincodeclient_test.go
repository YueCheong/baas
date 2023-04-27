package contract

import (
	"bit-bass/artifacts"
	"bit-bass/channel"
	"bit-bass/deploy"
	"bit-bass/network"
	"bit-bass/utils"
	"fmt"
	"testing"
	"time"
)

//此测试已经更新，启动的区块链网络将在artifact/262/目录中
const id = 262

func TestChaincodeClient(t *testing.T) {

	net := network.NewNodesNet("fabric_net")
	net.CreateNet()
	defer net.RemoveNet()

	c := deploy.NewConfigurator(utils.ConfigPathWithId(id), net)

	fmt.Println("set orderer org -> ", c.SetOrdererOrg("orderer", "example.com", "OrdererMSP"))
	fmt.Println("set orderer node -> ", c.SetOrdererNode("orderer", "7050/tcp"))

	fmt.Println("add peer org1 -> ", c.AddPeerOrg("Org1", "org1.example.com", "Org1MSP"))
	fmt.Println("add peer org2 -> ", c.AddPeerOrg("Org2", "org2.example.com", "Org2MSP"))

	fmt.Println("add peer0 node to org1 -> ", c.AddPeerNodeToOrg("Org1", "peer0", "7051/tcp"))
	fmt.Println("add peer1 node to org1 -> ", c.AddPeerNodeToOrg("Org1", "peer1", "8051/tcp"))
	fmt.Println("add peer0 node to org2 -> ", c.AddPeerNodeToOrg("Org2", "peer0", "9051/tcp"))
	fmt.Println("add peer1 node to org2 -> ", c.AddPeerNodeToOrg("Org2", "peer1", "10051/tcp"))

	fmt.Println("set cli connect peer0 node of org1 -> ", c.SetToolsCli("cli", c.Peerorgsconf[0].GetNodes()[0]))

	//Creat Aritifacts
	fmt.Println("Creat Artificats")
	fmt.Println("generate yaml - >", artifacts.GenerateYaml(c, utils.ConfigPathWithId(id)))
	fmt.Println("generate crypto - >", artifacts.GenerateCryptoConfig(utils.ConfigPathWithId(id)))
	fmt.Println("generate genesis - >", artifacts.GenerateGenesisBlock(utils.ConfigPathWithId(id)))

	n := network.NewNetManager()
	fmt.Println("load info -> ", n.LoadConfig(c))
	fmt.Println("prepared -> ", n.Prepared())
	fmt.Println("start network -> ", n.StartNet())

	var peerorg1, peerorg2 *network.PeerOrg
	for _, peerorg := range n.Peerorgs {
		if peerorg.Name() == "Org1" {
			peerorg1 = peerorg
		}
		if peerorg.Name() == "Org2" {
			peerorg2 = peerorg
		}
	}

	//Network started, start testing
	//Waiting for orderer container to be ready
	time.Sleep(time.Second)

	fmt.Println("ChannelClient test started:")

	//创建Channelmanager
	chm := channel.NewChannelmanager(utils.ConfigPathWithId(id))

	//在Channelmanager中建立新channelconfig
	chm.NewChannelInfo("mych")
	chi, ok := chm.GetChannelInfo("mych")
	if ok == false {
		//return nil, errors.New("Can't get channel info")
	}

	//创建sdk所需配置文件
	fmt.Println("Generate sdk info file -> ", channel.WriteSdkConfig(c, chm))

	var c1, c2 channel.ChannelIf

	//通过channelconfig建立操纵通道的channel对象
	c1, err := chi.NewChannelClient("Admin", peerorg1)
	fmt.Println("Creat channel object as org1-> ", err)

	c2, err = chi.NewChannelClient("Admin", peerorg2)
	fmt.Println("Creat channel object as org2-> ", err)

	//在完成对通道到操作后释放fabric sdk
	defer c1.Close()
	defer c2.Close()

	fmt.Println("channel channelname -> ", c1.ChannelName())

	//在区块连网络中创建通道
	fmt.Print("creat channel -> ")
	_, err = c1.CreateChannel()
	fmt.Println(err)

	//将peer节点加入通道
	fmt.Println("join channel:")
	//for _, node := range peerorg1.GetPeerNodes() {
	//	fmt.Printf("%v join channel -> %v\n", node.ContainerName(), c1.JoinChannel(node.ContainerName()))
	//}
	//测试组织节点不完全加入通道的解决方法
	pnode := peerorg1.GetPeerNodes()[0]
	fmt.Printf("%v join channel -> %v\n", pnode.ContainerName(), c1.JoinChannel(pnode.ContainerName()))

	for _, node := range peerorg2.GetPeerNodes() {
		fmt.Printf("%v join channel -> %v\n", node.ContainerName(), c2.JoinChannel(node.ContainerName()))
	}

	//通道和节点发生变化后需要更新sdk到参数文件，重新建立sdk Client
	fmt.Println("Regenerate sdk info after peer join channel -> ", channel.WriteSdkConfig(c, chm))
	fmt.Println("Refresh sdk with new info -> ", c1.RefreshSdk())

	//更新anchor节点
	fmt.Println("Update anchor peer:")
	fmt.Print("set anchor peer of org1 -> ")
	_, err = c1.UpdateAnchorPeers()
	fmt.Println(err)

	fmt.Print("set anchor peer of org2 -> ")
	_, err = c2.UpdateAnchorPeers()
	fmt.Println(err)

	//return n, nil

	//启动测试网络BYFN
	//
	//测试结束后关闭并清理网络
	defer func() {
		n.StopNet()
		n.Remove()
	}()

	//var peerorg1, peerorg2 *network.PeerOrg
	//peerorg1 = n.Peerorgs[0]
	//peerorg2 = n.Peerorgs[1]

	//创建chaincode
	ccConf := ChaincodeConfig{
		ChannelName:      "mych",
		ChaincodeName:    "mycc",
		ChaincodeDesc:    "An example chaincode to test fabric",
		ChaincodeGoPath:  utils.ConfigPath() + "/chaincode/",
		ChaincodePath:    "chaincode_example02/go/",
		ChaincodeVersion: "1.0",
	}
	//创建合约管理器
	cm := NewContractManager(utils.ConfigPathWithId(id))

	//建立新的合约
	ci, err := cm.NewChaincodeInfo(ccConf)
	fmt.Println("Creat Chaincode Info ->", err)

	//建立合约客户端
	fmt.Print("Creat Chaincode client -> ")
	cCli1, err := ci.NewChaincodeClient("Admin", peerorg1)
	if err != nil {
		fmt.Println(err)
	}
	cCli2, err := ci.NewChaincodeClient("Admin", peerorg2)
	fmt.Println(err)

	//打印合约信息
	fmt.Println("ChaincodeName -> ", cCli1.ChaincodeID())
	fmt.Println("Chaincode Version -> ", cCli1.ChaincodeVer())
	fmt.Println("Chaincode Desc -> ", cCli1.ChaincodeDesc())
	fmt.Println("Chaincode Channel -> ", cCli1.ChaincodeChan())

	//安装合约
	fmt.Print("Chaincode Install ->")
	installResp1, err := cCli1.InstallChaincode()
	if err != nil {
		fmt.Println(err)
	}
	installResp2, err := cCli2.InstallChaincode()
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(installResp1)
		fmt.Println(installResp2)
	}

	//实例化合约
	fmt.Print("Chaincode Init ->")
	_, err = cCli1.InstantiateChaincode("init", "a", "100", "b", "200")
	fmt.Println(err)

	time.Sleep(time.Second)

	fmt.Print("Chaincode Query ->")
	queryResp, err := cCli2.QueryChaincode("query", "a")
	if err != nil {
		fmt.Println("The error of query is :", err)
	} else {
		fmt.Println(queryResp.Responses[0].Response)
	}

	fmt.Print("Chaincode Invoke ->")
	invokeResp, err := cCli1.InvokeChaincode("invoke", "a", "b", "30")
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(invokeResp.Responses[0].Response)
	}

	fmt.Print("Chaincode Query ->")
	queryResp, err = cCli1.QueryChaincode("query", "a")
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(queryResp.Responses[0].Response)
	}

	//进行合约升级测试
	fmt.Println("Chaincode upgrade test :")
	ccConf.ChaincodeVersion = "1.1"
	newci, _ := cm.NewChaincodeInfo(ccConf)

	//创建新版本合约客户端
	fmt.Print("Creat new cc client ->")
	newcCli1, err := newci.NewChaincodeClient("Admin", peerorg1)
	if err != nil {
		fmt.Println(err)
	}
	newcCli2, err := newci.NewChaincodeClient("Admin", peerorg2)
	fmt.Println(err)

	//安装新版合约
	fmt.Print("New Chaincode Install ->")
	installResp1, err = newcCli1.InstallChaincode()
	if err != nil {
		fmt.Println(err)
	}
	installResp2, err = newcCli2.InstallChaincode()
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(installResp1)
		fmt.Println(installResp2)
	}

	//升级合约
	fmt.Print("Chaincode upgrade ->")
	_, err = newcCli1.UpgradeChaincode("init", "a", "100", "b", "200")
	fmt.Println(err)

	time.Sleep(time.Second)

	//测试升级后的合约
	fmt.Println("Test new version of contract")

	fmt.Print("Chaincode Query ->")
	queryResp, err = newcCli2.QueryChaincode("query", "a")
	if err != nil {
		fmt.Println("The error of query is :", err)
	} else {
		fmt.Println(queryResp.Responses[0].Response)
	}

	fmt.Print("Chaincode Invoke ->")
	invokeResp, err = newcCli1.InvokeChaincode("invoke", "a", "b", "30")
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(invokeResp.Responses[0].Response)
	}

	time.Sleep(time.Second)

	fmt.Print("Chaincode Query ->")
	queryResp, err = newcCli1.QueryChaincode("query", "a")
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(queryResp.Responses[0].Response)
	}

	fmt.Println("Test channel client, get contract:")
	ccResp, err := c1.GetContracts(peerorg1.GetPeerNodes()[0])
	fmt.Println("c1 get contract ->", ccResp.Chaincodes, err)

	fmt.Println("Chaincode test complete")
}

//
//func StartupBYFN_Network() (*network.NetManager, error) {
//
//	c := deploy.NewConfigurator(utils.ConfigPathWithId(id), "fabric_net")
//	fmt.Println("set orderer org -> ", c.SetOrdererOrg("orderer", "example.com", "OrdererMSP"))
//	fmt.Println("set orderer node -> ", c.SetOrdererNode("orderer", "7050/tcp"))
//
//	fmt.Println("add peer org1 -> ", c.AddPeerOrg("Org1", "org1.example.com", "Org1MSP"))
//	fmt.Println("add peer org2 -> ", c.AddPeerOrg("Org2", "org2.example.com", "Org2MSP"))
//
//	fmt.Println("add peer0 node to org1 -> ", c.AddPeerNodeToOrg("Org1", "peer0", "7051/tcp"))
//	fmt.Println("add peer1 node to org1 -> ", c.AddPeerNodeToOrg("Org1", "peer1", "8051/tcp"))
//	fmt.Println("add peer0 node to org2 -> ", c.AddPeerNodeToOrg("Org2", "peer0", "9051/tcp"))
//	fmt.Println("add peer1 node to org2 -> ", c.AddPeerNodeToOrg("Org2", "peer1", "10051/tcp"))
//
//	fmt.Println("set cli connect peer0 node of org1 -> ", c.SetToolsCli("cli", c.Peerorgsconf[0].GetNodes()[0]))
//
//	//Creat Aritifacts
//	fmt.Println("Creat Artificats")
//	fmt.Println("generate yaml - >", artifacts.GenerateYaml(c, utils.ConfigPathWithId(id)))
//	fmt.Println("generate crypto - >", artifacts.GenerateCryptoConfig(utils.ConfigPathWithId(id)))
//	fmt.Println("generate genesis - >", artifacts.GenerateGenesisBlock(utils.ConfigPathWithId(id)))
//
//	n := network.NewNetManager()
//	fmt.Println("load info -> ", n.LoadConfig(c))
//	fmt.Println("prepared -> ", n.Prepared())
//	fmt.Println("start network -> ", n.StartNet())
//
//	var peerorg1, peerorg2 *network.PeerOrg
//	for _, peerorg := range n.Peerorgs {
//		if peerorg.Name() == "Org1" {
//			peerorg1 = peerorg
//		}
//		if peerorg.Name() == "Org2" {
//			peerorg2 = peerorg
//		}
//	}
//
//	//Network started, start testing
//	//Waiting for orderer container to be ready
//	time.Sleep(time.Second)
//
//	fmt.Println("ChannelClient test started:")
//
//	//创建Channelmanager
//	cm := channel.NewChannelmanager(utils.ConfigPathWithId(id))
//
//	//在Channelmanager中建立新channelconfig
//	cm.NewChannelInfo("mych")
//	ci, ok := cm.GetChannelInfo("mych")
//	if ok == false {
//		return nil, errors.New("Can't get channel info")
//	}
//
//	//创建sdk所需配置文件
//	fmt.Println("Generate sdk info file -> ", channel.WriteSdkConfig(c, cm))
//
//	var c1, c2 channel.ChannelIf
//
//	//通过channelconfig建立操纵通道的channel对象
//	c1, err := ci.NewChannelClient("Admin", peerorg1)
//	fmt.Println("Creat channel object as org1-> ", err)
//
//	c2, err = ci.NewChannelClient("Admin", peerorg2)
//	fmt.Println("Creat channel object as org2-> ", err)
//
//	//在完成对通道到操作后释放fabric sdk
//	defer c1.Close()
//	defer c2.Close()
//
//	fmt.Println("channel channelname -> ", c1.ChannelName())
//
//	//在区块连网络中创建通道
//	fmt.Print("creat channel -> ")
//	_, err = c1.CreateChannel()
//	fmt.Println(err)
//
//	//将peer节点加入通道
//	fmt.Println("join channel:")
//	//for _, node := range peerorg1.GetPeerNodes() {
//	//	fmt.Printf("%v join channel -> %v\n", node.ContainerName(), c1.JoinChannel(node.ContainerName()))
//	//}
//	//测试组织节点不完全加入通道的解决方法
//	pnode := peerorg1.GetPeerNodes()[0]
//	fmt.Printf("%v join channel -> %v\n", pnode.ContainerName(), c1.JoinChannel(pnode.ContainerName()))
//
//	for _, node := range peerorg2.GetPeerNodes() {
//		fmt.Printf("%v join channel -> %v\n", node.ContainerName(), c2.JoinChannel(node.ContainerName()))
//	}
//
//	//通道和节点发生变化后需要更新sdk到参数文件，重新建立sdk Client
//	fmt.Println("Regenerate sdk info after peer join channel -> ", channel.WriteSdkConfig(c, cm))
//	fmt.Println("Refresh sdk with new info -> ", c1.RefreshSdk())
//
//	//更新anchor节点
//	fmt.Println("Update anchor peer:")
//	fmt.Print("set anchor peer of org1 -> ")
//	_, err = c1.UpdateAnchorPeers()
//	fmt.Println(err)
//
//	fmt.Print("set anchor peer of org2 -> ")
//	_, err = c2.UpdateAnchorPeers()
//	fmt.Println(err)
//
//	return n, nil
//}
