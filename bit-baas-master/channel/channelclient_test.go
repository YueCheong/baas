package channel

import (
	"bit-bass/artifacts"
	"bit-bass/deploy"
	"bit-bass/network"
	"bit-bass/utils"
	"fmt"
	"testing"
	"time"
)

func TestChannel(t *testing.T) {

	//此测试已经更新，创建的区块链网络所需的artifacts将在artifacts/261目录下
	id := 261

	net := network.NewNodesNet("fabric_net")
	net.CreateNet()
	defer net.RemoveNet()

	c := deploy.NewConfigurator(utils.ConfigPathWithId(id), net)
	fmt.Println("set orderer org -> ", c.SetOrdererOrg("orderer", "example.com", "OrdererMSP"))
	fmt.Println("set orderer node -> ", c.SetOrdererNode("orderer", "5050/tcp"))

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

	//Shutdown network after testing
	defer func() {
		fmt.Println("stop network -> ", n.StopNet())
		fmt.Println("remove -> ", n.Remove())
	}()

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
	cm := NewChannelmanager(utils.ConfigPathWithId(id))

	//在Channelmanager中建立新channelconfig
	cm.NewChannelInfo("mych")
	ci, ok := cm.GetChannelInfo("mych")
	if ok == false {
		return
	}

	//创建sdk所需配置文件
	fmt.Println("Generate sdk info file -> ", WriteSdkConfig(c, cm))

	var c1, c2 ChannelIf

	//通过channelconfig建立操纵通道的channel对象
	c1, err := ci.NewChannelClient("Admin", peerorg1)
	fmt.Println("Creat channel object as org1-> ", err)

	c2, err = ci.NewChannelClient("Admin", peerorg2)
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
	for _, node := range peerorg1.GetPeerNodes() {
		fmt.Printf("%v join channel -> %v\n", node.ContainerName(), c1.JoinChannel(node.ContainerName()))
	}
	for _, node := range peerorg2.GetPeerNodes() {
		fmt.Printf("%v join channel -> %v\n", node.ContainerName(), c2.JoinChannel(node.ContainerName()))
	}

	//通道和节点发生变化后需要更新sdk到参数文件，重新建立sdk Client
	fmt.Println("Regenerate sdk info after peer join channel -> ", WriteSdkConfig(c, cm))
	fmt.Println("Refresh sdk with new info -> ", c1.RefreshSdk())

	//更新anchor节点
	fmt.Println("Update anchor peer:")
	fmt.Print("set anchor peer of org1 -> ")
	_, err = c1.UpdateAnchorPeers()
	fmt.Println(err)

	fmt.Print("set anchor peer of org2 -> ")
	_, err = c2.UpdateAnchorPeers()
	fmt.Println(err)

	fmt.Print("get blockchain info -> ")
	blockchain, err := c1.GetBlockchainInfo()
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(blockchain.BCI)
	}

	fmt.Print("query block by height -> ")
	block1, err := c1.GetBlockByHeight(0)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(block1.Header)
	}

	fmt.Print("query blcok by hash -> ")
	block2, err := c1.GetBlcokByHash(blockchain.BCI.CurrentBlockHash)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(block2.Header)
	}

	//没有交易ID，因此这个函数暂时无法测试
	//fmt.Print("query transaction by txid ->")
	//tx,err := c1.GetTransaction(string(creatchannelTxId))
	//if err != nil {
	//	fmt.Println(err)
	//} else {
	//	fmt.Println(tx)
	//}
	ccResp, err := c1.GetContracts(peerorg1.GetPeerNodes()[0])
	if err != nil {
		fmt.Println("Get contract enconter error")
	}
	fmt.Println("get contracts -> ", ccResp, err)

	//返回某个节点加入了那些通道
	fmt.Print("get channels that peer0.org1 joined ->")
	chs, err := c1.QueryChannelofPeer(peerorg1.GetPeerNodes()[0])
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(chs)
	}

	//需要注意，只有对应组织的channelclient才能查询对应组织内的节点。
	//即下面代码中的c2的身份为org2，能够查询org2中的peer加入了那些通道
	//如果创建的channelclient身份为org1，则不能查询org中的peer节点加入
	//了哪些通道，否则会返回access denied错误
	fmt.Print("get channels that peer0.org2 joined ->")
	chs, err = c2.QueryChannelofPeer(peerorg2.GetPeerNodes()[0])
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(chs)
	}

	fmt.Println("ChannelClient type test complete")

}
