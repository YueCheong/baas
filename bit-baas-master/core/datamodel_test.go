package core

//
//import (
//	"bit-bass/contract"
//	"encoding/json"
//	"fmt"
//	"reflect"
//	"testing"
//)
//
//func TestNetworkToJson(t *testing.T) {
//	orign := BlockchainToJson{
//		ID:   1,
//		Name: "firstnet",
//		OrdererOrg: OrdererOrgToJson{
//			Name:   "Orderer",
//			Domain: "Testing.com",
//			MSPID:  "OrdererMSP",
//			Orderer: NodeToJson{
//				Name: "Orderer",
//				Port: "5050",
//			},
//		},
//		PeerOrg: []PeerOrgToJson{PeerOrgToJson{
//			Name:   "org1",
//			Domain: "Testing.com",
//			MSPID:  "Org1MSP",
//			Peers: []NodeToJson{
//				NodeToJson{
//					Name: "peer0",
//					Port: "5051",
//				},
//				{
//					Name: "peer1",
//					Port: "6051",
//				},
//			},
//		}},
//		Config: nil,
//		Channels: []ChannelToJson{
//			{
//				Name:  "mychannel",
//				Peers: []string{"peer0.org1.Testing.com"},
//			},
//		},
//	}
//
//	networkJson, err := json.Marshal(orign)
//	if err != nil {
//		panic("Marshal Failed")
//	}
//	fmt.Println(string(networkJson))
//
//	var fromJson BlockchainToJson
//	err = json.Unmarshal(networkJson, &fromJson)
//	if err != nil {
//		panic("Unmarshal Failed")
//	}
//
//	fmt.Println(fromJson)
//	fmt.Println("---------------------------------")
//
//	if !reflect.DeepEqual(orign, fromJson) {
//		panic("Data Error")
//	}
//}
//
//func TestChannelManageInfoToJson(t *testing.T) {
//	orign := ChannelManageInfoToJson{
//		ChannelName: "mychannel",
//		Operation:   AddPeerToChannel,
//		Args:        []string{"peer1.org1.Testing.com"},
//	}
//
//	dataJson, err := json.Marshal(orign)
//	if err != nil {
//		panic("Marshal Failed")
//	}
//	fmt.Println(string(dataJson))
//
//	var fromJson ChannelManageInfoToJson
//	err = json.Unmarshal(dataJson, &fromJson)
//	if err != nil {
//		panic("Unmarshal Failed")
//	}
//
//	fmt.Println(fromJson)
//	fmt.Println("---------------------------------")
//
//	if !reflect.DeepEqual(orign, fromJson) {
//		panic("Data Error")
//	}
//}
//
//func TestContractTOJson(T *testing.T) {
//	orign := ContractToJson{
//		ID:              1,
//		ContractName:    "simplecc",
//		ChannelName:     "mychannel",
//		ContractDesc:    "A simple example of chaincode",
//		ContractVersion: "1.0",
//		ContractPath:    "/chaincode_example02/go/",
//		ContractLang:    contract.Golang,
//	}
//
//	dataJson, err := json.Marshal(orign)
//	if err != nil {
//		panic("Marshal Failed")
//	}
//	fmt.Println(string(dataJson))
//
//	var fromJson ContractToJson
//	err = json.Unmarshal(dataJson, &fromJson)
//	if err != nil {
//		panic("Unmarshal Failed")
//	}
//
//	fmt.Println(fromJson)
//	fmt.Println("---------------------------------")
//
//	if !reflect.DeepEqual(orign, fromJson) {
//		panic("Data Error")
//	}
//}
//
//func TestContractManageInfoToJson(T *testing.T) {
//	orign := ContractManageInfoToJson{
//		ID:              1,
//		BlockchainID:    1,
//		ContractName:    "mycc",
//		ContractVersion: "1.0",
//		Operation:       Instantiate,
//		Args:            []string{"init", "a", "100", "b", "100"},
//	}
//
//	dataJson, err := json.Marshal(orign)
//	if err != nil {
//		panic("Marshal Failed")
//	}
//	fmt.Println(string(dataJson))
//
//	var fromJson ContractManageInfoToJson
//	err = json.Unmarshal(dataJson, &fromJson)
//	if err != nil {
//		panic("Unmarshal Failed")
//	}
//
//	fmt.Println(fromJson)
//	fmt.Println("---------------------------------")
//
//	if !reflect.DeepEqual(orign, fromJson) {
//		panic("Data Error")
//	}
//}
//
//func TestContractInvokeInfoToJson(T *testing.T) {
//	orign := ContractInvokeInfoToJson{
//		ID:              1,
//		BlockchainID:    1,
//		ContractName:    "mycc",
//		ContractVersion: "1.0",
//		InvokeType:      Invoke,
//		Args:            []string{"transfer", "a", "b", "50"},
//	}
//
//	dataJson, err := json.Marshal(orign)
//	if err != nil {
//		panic("Marshal Failed")
//	}
//	fmt.Println(string(dataJson))
//
//	var fromJson ContractInvokeInfoToJson
//	err = json.Unmarshal(dataJson, &fromJson)
//	if err != nil {
//		panic("Unmarshal Failed")
//	}
//
//	fmt.Println(fromJson)
//	fmt.Println("---------------------------------")
//
//	if !reflect.DeepEqual(orign, fromJson) {
//		panic("Data Error")
//	}
//}
