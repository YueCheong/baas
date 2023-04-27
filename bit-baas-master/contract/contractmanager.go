package contract

import (
	"bit-bass/utils"
	"errors"
)

type ContractManager struct {
	configPath string
	idGen      *utils.AutoIncIDGen
	chaincodes map[int]*ChaincodeInfo
}

func NewContractManager(configPath string) *ContractManager {
	cm := ContractManager{}
	cm.idGen = utils.NewAutoIncID()
	cm.chaincodes = make(map[int]*ChaincodeInfo)
	cm.configPath = configPath
	return &cm
}

func (cm *ContractManager) NewChaincodeInfo(c ChaincodeConfig) (*ChaincodeInfo, error) {
	for i, _ := range cm.chaincodes {
		if cm.chaincodes[i].chaincodeName == c.ChaincodeName &&
			cm.chaincodes[i].channelName == c.ChannelName &&
			cm.chaincodes[i].chaincodeVersion == c.ChaincodeVersion {
			return nil, errors.New("Duplicate Chaincode")
		}
	}

	ci := ChaincodeInfo{
		id:               cm.idGen.GenID(),
		configPath:       cm.configPath,
		channelName:      c.ChannelName,
		chaincodeName:    c.ChaincodeName,
		chaincodeDesc:    c.ChaincodeDesc,
		chaincodeGoPath:  c.ChaincodeGoPath,
		chaincodePath:    c.ChaincodePath,
		chaincodeVersion: c.ChaincodeVersion,
		chaincodeLang:    c.ChaincodeLang,
	}

	cm.chaincodes[ci.id] = &ci
	return &ci, nil
}

func (cm *ContractManager) GetChaincodes() []*ChaincodeInfo {
	var outputCi []*ChaincodeInfo
	for i := range cm.chaincodes {
		outputCi = append(outputCi, cm.chaincodes[i])
	}
	return outputCi
}

func (cm *ContractManager) GetChaincodeByID(id int) (*ChaincodeInfo, bool) {
	ci, ok := cm.chaincodes[id]
	return ci, ok
}

func (cm *ContractManager) GetChaincodeByName(ccname string) []*ChaincodeInfo {
	var outputCi []*ChaincodeInfo
	for i, _ := range cm.chaincodes {
		if cm.chaincodes[i].chaincodeName == ccname {
			outputCi = append(outputCi, cm.chaincodes[i])
		}
	}
	return outputCi
}

func (cm *ContractManager) GetChaincodeByChannelId(chid string) []*ChaincodeInfo {
	var outputCi []*ChaincodeInfo
	for i, _ := range cm.chaincodes {
		if cm.chaincodes[i].channelName == chid {
			outputCi = append(outputCi, cm.chaincodes[i])
		}
	}
	return outputCi
}
