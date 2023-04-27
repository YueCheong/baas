package core

import (
	"bit-bass/artifacts"
	"bit-bass/channel"
	"bit-bass/deploy"
	"bit-bass/network"
	"bit-bass/utils"
	"fmt"
	"github.com/pkg/errors"
	"testing"
	"time"
)

func TestBlockchainManager(t *testing.T) {
	bm := NewBlockchainManager()

	byfn := BlockchainToJson{
		Name: "firstnet",
		OrdererOrg: OrdererOrgToJson{
			Name:   "Orderer",
			Domain: "testing.com",
			MSPID:  "OrdererMSP",
			Orderer: NodeToJson{
				Name: "Orderer",
				Port: "7050",
			},
		},
		PeerOrg: []PeerOrgToJson{{
			Name:   "org1",
			Domain: "testing.com",
			MSPID:  "Org1MSP",
			Peers: []NodeToJson{
				{
					Name: "peer0",
					Port: "7051",
				},
				{
					Name: "peer1",
					Port: "8051",
				},
			},
		}, {
			Name:   "org2",
			Domain: "testing.com",
			MSPID:  "Org2MSP",
			Peers: []NodeToJson{
				{
					Name: "peer0",
					Port: "9051",
				}, {
					Name: "peer1",
					Port: "10051",
				},
			},
		},
		},
		Config: nil,
		Channels: []ChannelToJson{
			{
				Name: "mych",
				Peers: []string{"peer0.org1.testing.com", "peer1.org1.testing.com",
					"peer0.org2.testing.com", "peer1.org2.testing.com"},
			},
		},
	}

	id := 259
	id, err := bm.NewBlockchain(byfn)
	fmt.Println("Creat new blockchain ->")
	if err != nil {
		t.Fail()
		fmt.Println(err)
	}

	_, ok := bm.GetBlockchainById(id)
	if !ok {
		panic("Can't get blockchain")
	}

	err = bm.StartBlockchainById(id)
	if err != nil {
		fmt.Println("Failed to start blockchain : ", err)
	}

	fmt.Print("Waiting ->")
	time.Sleep(time.Second * 2)

	fmt.Print("Stop and remove Blockchain")
	err = bm.StopAndRemoveBlockchainById(id)
	if err != nil {
		t.Fail()
		fmt.Println(err)
	}

	for _, net := range bm.GetNets() {
		bm.RemoveNetByName(net.NetName())
	}

}

func StartupBYFN_Network(id int) (*network.NetManager, error) {

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
	cm := channel.NewChannelmanager(utils.ConfigPathWithId(id))

	//在Channelmanager中建立新channelconfig
	cm.NewChannelInfo("mych")
	ci, ok := cm.GetChannelInfo("mych")
	if ok == false {
		return nil, errors.New("Can't get channel info")
	}

	//创建sdk所需配置文件
	fmt.Println("Generate sdk info file -> ", channel.WriteSdkConfig(c, cm))

	var c1, c2 channel.ChannelIf

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
	fmt.Println("Regenerate sdk info after peer join channel -> ", channel.WriteSdkConfig(c, cm))
	fmt.Println("Refresh sdk with new info -> ", c1.RefreshSdk())

	//更新anchor节点
	fmt.Println("Update anchor peer:")
	fmt.Print("set anchor peer of org1 -> ")
	_, err = c1.UpdateAnchorPeers()
	fmt.Println(err)

	fmt.Print("set anchor peer of org2 -> ")
	_, err = c2.UpdateAnchorPeers()
	fmt.Println(err)

	return n, nil
}
