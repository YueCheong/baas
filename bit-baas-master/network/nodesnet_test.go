package network

import (
	"fmt"
	"testing"
)

func TestNodesNet(t *testing.T) {
	fmt.Println("Start testing DockConnter")

	net := NewDefaultNodesNet(FABRICNET)

	fmt.Println("net name is -> ", net.NetName())

	fmt.Println("create net return -> ", net.CreateNet())

	fmt.Println("net id is -> ", net.NetID())

	fmt.Println("test net exist -> ", net.IsNetExist())

	fmt.Println("net info -> ", net.Inspect())

	fmt.Println("remove net return -> ", net.RemoveNet())

	fmt.Println("test net exist -> ", net.IsNetExist())
}
