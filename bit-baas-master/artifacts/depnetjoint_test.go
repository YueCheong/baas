package artifacts

import (
	"bit-bass/deploy"
	"bit-bass/network"
	"bit-bass/utils"
	"fmt"
	"testing"
	"time"
)

func TestJoint(t *testing.T) {
	net := network.NewNodesNet("fabric_net")
	net.CreateNet()
	c := *deploy.NewConfigurator(utils.ConfigPathWithId(id), net)
	fmt.Println("set orderer org -> ", c.SetOrdererOrg("orderer", "example.com", "OrdererMSP"))
	fmt.Println("set orderer node -> ", c.SetOrdererNode("orderer", "7050/tcp"))
	fmt.Println("add peer org1 -> ", c.AddPeerOrg("Org1", "org1.example.com", "Org1MSP"))
	fmt.Println("add peer org2 -> ", c.AddPeerOrg("Org2", "org2.example.com", "Org2MSP"))
	fmt.Println("add peer0 node to org1 -> ", c.AddPeerNodeToOrg("Org1", "peer0", "7051/tcp"))
	fmt.Println("add peer1 node to org1 -> ", c.AddPeerNodeToOrg("Org1", "peer1", "8051/tcp"))
	fmt.Println("add peer0 node to org2 -> ", c.AddPeerNodeToOrg("Org2", "peer0", "9051/tcp"))
	fmt.Println("add peer1 node to org2 -> ", c.AddPeerNodeToOrg("Org2", "peer1", "10051/tcp"))

	fmt.Println("set cli connect peer0 node of org1 -> ", c.SetToolsCli("cli", c.Peerorgsconf[0].GetNodes()[0]))

	//deploy
	fmt.Println("generate config -> ", GenerateYaml(&c, utils.ConfigPathWithId(id)))
	fmt.Println("generate crypto -> ", GenerateCryptoConfig(utils.ConfigPathWithId(id)))
	fmt.Println("generte genesis-> ", GenerateGenesisBlock(utils.ConfigPathWithId(id)))

	//network
	n := network.NewNetManager()
	fmt.Println("load config -> ", n.LoadConfig(&c))
	fmt.Println("prepared -> ", n.Prepared())
	fmt.Println("start network -> ", n.StartNet())
	time.Sleep(10 * time.Second)
	fmt.Println(n.Peerorgs[0].GetPeerNodes()[0].PrintLog())
	fmt.Println("stop network -> ", n.StopNet())
	fmt.Println("remove -> ", n.Remove())
}
