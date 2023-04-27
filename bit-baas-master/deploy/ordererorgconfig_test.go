package deploy

import (
	"bit-bass/utils"
	"fmt"
	"testing"
)

func TestNewOrdererOrgConfig(t *testing.T) {
	path := utils.ConfigPath()
	fmt.Println(path)
	ordererorg := NewOrdererOrgConfig("orderer", "example.com", "OrdererMSP", path)
	fmt.Println(ordererorg)
}
