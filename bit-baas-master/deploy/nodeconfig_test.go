package deploy

import (
	"bit-bass/utils"
	"fmt"
	"testing"
)

func TestNodeConfig(t *testing.T) {
	node := NewNodeConfig("peer0", nil, "tcp/7051", "fabric_net")
	fmt.Println(node)
}

func TestPeerOrg(t *testing.T) {
	path := utils.ConfigPath()
	peerorg := NewPeerOrgConfig("test", "org1.example.com", "Org1MSP", path)
	peer := NewNodeConfig("peer0", nil, "tcp/7051", "fabric_net")
	if err := peer.JoinOrg(peerorg); err != nil {
		panic(err)
	}
	fmt.Println(peerorg)
	fmt.Println(peer)
}

func TestOrdererOrg(t *testing.T) {
	path := utils.ConfigPath()
	ordererorg := NewOrdererOrgConfig("orderer", "example.com", "OrdererMSP", path)
	orderer := NewNodeConfig("orderer", nil, "tcp/7050", "fabric_net")
	if err := orderer.JoinOrg(ordererorg); err != nil {
		panic(err)
	}
	fmt.Println(ordererorg)
	fmt.Println(orderer)
}
