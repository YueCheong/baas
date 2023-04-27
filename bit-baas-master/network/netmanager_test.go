package network

import (
	"bit-bass/artifacts"
	"bit-bass/deploy"
	"bit-bass/utils"
	"fmt"
	"testing"
	"time"
)

//
const id = 263

func TestConfigLoad(t *testing.T) {
	net := NewNodesNet("fabric_net")
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

	//deploy
	//fmt.Println("generate config -> ", artifacts.GenerateConfig(*c))
	//fmt.Println("execute config -> ", artifacts.ExecuteBat())

	n := NewNetManager()
	fmt.Println("load config -> ", n.LoadConfig(c))
	fmt.Println("prepared -> ", n.Prepared())
	fmt.Println("start network -> ", n.StartNet())
	//time.Sleep(2 * time.Second)
	fmt.Println(n.Ordererorg.orderer.PrintLog())

	fmt.Println(">>>>> here make a break point <<<<<")
	fmt.Println()

	fmt.Println("stop network -> ", n.StopNet())
	fmt.Println("remove -> ", n.Remove())

}

func TestStartup(t *testing.T) {
	path := utils.ConfigPathWithId(id)
	net := NewDefaultNodesNet(FABRICNET)
	fmt.Println("create net return -> ", net.CreateNet())

	ordererconfig := deploy.NewOrdererOrgConfig("orderer", "example.com", "OrdererMSP", path)
	deploy.NewNodeConfig("orderer", ordererconfig, "7050/tcp", FABRICNET)
	ordererorg := NewOrdererOrg(ordererconfig)

	org1config := deploy.NewPeerOrgConfig("Org1", "org1.example.com", "Org1MSP", path)
	deploy.NewNodeConfig("peer0", org1config, "7051/tcp", FABRICNET)
	deploy.NewNodeConfig("peer1", org1config, "8051/tcp", FABRICNET)
	peerorg1 := NewPeerOrg(org1config)

	org2config := deploy.NewPeerOrgConfig("Org2", "org2.example.com", "Org2MSP", path)
	deploy.NewNodeConfig("peer0", org2config, "9051/tcp", FABRICNET)
	deploy.NewNodeConfig("peer1", org2config, "10051/tcp", FABRICNET)
	peerorg2 := NewPeerOrg(org2config)

	onode := ordererorg.GetOrdererNode()
	_ = onode.AddOrgCrypto(org1config)
	_ = onode.AddOrgCrypto(org2config)

	fmt.Println("create ordererorg -> ", ordererorg.Create())
	fmt.Println("create peerorg1 -> ", peerorg1.Create())
	fmt.Println("create peerorg2 -> ", peerorg2.Create())

	fmt.Println("ordererorg connect to net -> ", ordererorg.ConnectNet(net))
	fmt.Println("peerorg1 connect to net -> ", peerorg1.ConnectNet(net))
	fmt.Println("peerorg2 connect to net -> ", peerorg2.ConnectNet(net))

	fmt.Println("start ordererorg -> ", ordererorg.Start())
	fmt.Println("start peerorg1 -> ", peerorg1.Start())
	fmt.Println("start peerorg2 -> ", peerorg2.Start())

	time.Sleep(10 * time.Second)

	fmt.Println("stop ordererorg -> ", ordererorg.Stop())
	fmt.Println("stop peerorg1 -> ", peerorg1.Stop())
	fmt.Println("stop peerorg2 -> ", peerorg2.Stop())

	fmt.Println("remove ordererorg -> ", ordererorg.Remove())
	fmt.Println("remove peerorg1 -> ", peerorg1.Remove())
	fmt.Println("remove peerorg2 -> ", peerorg2.Remove())
	fmt.Println("remove net return -> ", net.RemoveNet())

}
