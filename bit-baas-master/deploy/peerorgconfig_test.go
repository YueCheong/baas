package deploy

import (
	"bit-bass/utils"
	"fmt"
	"testing"
)

func TestNewPeerOrgConfig(t *testing.T) {
	path := utils.ConfigPath()
	fmt.Println(path)
	peerorg := NewPeerOrgConfig("test", "org1.example.com", "Org1MSP", path)
	fmt.Println(peerorg)
}
