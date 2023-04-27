package network

import (
	"bit-bass/deploy"
	"bit-bass/utils"
	"fmt"
	"testing"
	"time"
)

func TestOrdererOrg(t *testing.T) {
	net := NewDefaultNodesNet(FABRICNET)
	fmt.Println("create net return -> ", net.CreateNet())

	path := utils.ConfigPathWithId(id)
	ordererconfig := deploy.NewOrdererOrgConfig("orderer", "example.com", "OrdererMSP", path)
	deploy.NewNodeConfig("orderer", ordererconfig, "7050/tcp", FABRICNET)

	ordererorg := NewOrdererOrg(ordererconfig)

	fmt.Println("create all -> ", ordererorg.Create())
	fmt.Println("connect to net -> ", ordererorg.ConnectNet(net))
	fmt.Println("start all -> ", ordererorg.Start())

	node := ordererorg.GetOrdererNode()
	fmt.Println(node)

	time.Sleep(10 * time.Second)
	fmt.Println(node.PrintLog())

	fmt.Println("stop all -> ", ordererorg.Stop())
	fmt.Println("remove all -> ", ordererorg.Remove())

	fmt.Println("remove net return -> ", net.RemoveNet())

}
