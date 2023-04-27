package network

import (
	"bit-bass/deploy"
	"bit-bass/utils"
	"fmt"
	"testing"
	"time"
)

func TestPeerOrg(t *testing.T) {
	net := NewDefaultNodesNet(FABRICNET)
	fmt.Println("create net return -> ", net.CreateNet())

	path := utils.ConfigPathWithId(id)
	org1config := deploy.NewPeerOrgConfig("test", "org1.example.com", "Org1MSP", path)
	deploy.NewNodeConfig("peer0", org1config, "7051/tcp", FABRICNET)
	deploy.NewNodeConfig("peer1", org1config, "8051/tcp", FABRICNET)

	peerorg := NewPeerOrg(org1config)

	fmt.Println("create all -> ", peerorg.Create())
	fmt.Println("connect to net -> ", peerorg.ConnectNet(net))
	fmt.Println("start all -> ", peerorg.Start())

	node1 := peerorg.GetPeerNodes()[0]
	node2 := peerorg.GetPeerNodes()[1]
	fmt.Println(node1)
	fmt.Println(node2)
	time.Sleep(10 * time.Second)
	fmt.Println(node1.PrintLog())

	fmt.Println("stop all -> ", peerorg.Stop())
	fmt.Println("remove all -> ", peerorg.Remove())

	fmt.Println("remove net return -> ", net.RemoveNet())

}
