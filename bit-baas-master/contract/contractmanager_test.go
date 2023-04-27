package contract

import (
	"bit-bass/utils"
	"fmt"
	"testing"
)

func TestContractManager(T *testing.T) {
	fmt.Println("Test ContractManager")

	m := NewContractManager(utils.ConfigPath())

	ccConfig := ChaincodeConfig{
		ChannelName:      "mychannel",
		ChaincodeName:    "mycc",
		ChaincodeDesc:    "An example chaincode to test fabric",
		ChaincodeGoPath:  utils.ConfigPath() + "/chaincode/",
		ChaincodePath:    "chaincode_example02/go/",
		ChaincodeVersion: "1.0",
	}

	ccConfig.ChaincodeVersion = "1.1"
	_, err := m.NewChaincodeInfo(ccConfig)
	fmt.Println("Creat New chaincode with new version -> ", err)

	ccConfig.ChaincodeName = "newcc"
	_, err = m.NewChaincodeInfo(ccConfig)

	fmt.Println("Get chaincodes -> ")
	for _, cc := range m.GetChaincodes() {
		fmt.Println(cc)
	}

	fmt.Println("Get chaincodes with id 'mycc' -> ")
	for _, cc := range m.GetChaincodeByName("mycc") {
		fmt.Println(cc)
	}

	fmt.Println("Get chaincodes with channel 'newchannel' -> ")
	for _, cc := range m.GetChaincodeByChannelId("newchannel") {
		fmt.Println(cc)
	}

}
