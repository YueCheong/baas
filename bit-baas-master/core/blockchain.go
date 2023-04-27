package core

import (
	"bit-bass/artifacts"
	"bit-bass/channel"
	"bit-bass/contract"
	"bit-bass/deploy"
	"bit-bass/logger"
	"bit-bass/network"
	"bit-bass/utils"
	"fmt"
	"github.com/pkg/errors"
	"os"
)

type Blockchain struct {
	id              int
	name            string
	netManager      *network.NetManager
	channelManager  *channel.ChannelManager
	contractManager *contract.ContractManager
	configurator    *deploy.Configurator
	ccLogger        *logger.ContractInvokeLogger
	status          BlockchainStatus
	dockerNetName   string
}

//获取Blockchain的信息
func (b *Blockchain) GetId() int {
	return b.id
}

func (b *Blockchain) GetName() string {
	return b.name
}

func (b *Blockchain) GetStatus() BlockchainStatus {
	return b.status
}

func (b *Blockchain) GetDockerNetName() string {
	return b.dockerNetName
}

//获取Blockchain中的Manager
func (b *Blockchain) GetChannelManager() *channel.ChannelManager {
	return b.channelManager
}

func (b *Blockchain) GetContractManager() *contract.ContractManager {
	return b.contractManager
}

func (b *Blockchain) GetConfigurator() *deploy.Configurator {
	return b.configurator
}

//Channel相关
func (b *Blockchain) CreatAndInitChannel(chName string, peers []string, anchors map[string]string) error {
	//由于对属于不同组织的peer进行操作需要以不同组织权限建立的channel client
	//创建一个临时Map，用于存储在建立和设置过程中创建的client
	//避免反复关闭和建立client
	var chClients = make(map[string]*channel.ChannelClient)

	if len(peers) == 0 {
		return errors.New("Can't creat channel without peer joined")
	}

	//建立通道信息
	ci, err := b.channelManager.NewChannelInfo(chName)
	if err != nil {
		return errors.Errorf("Failed to creat channel : %v", err)
	}

	//创建channel client所需的sdk config
	err = channel.WriteSdkConfig(b.configurator, b.channelManager)
	if err != nil {
		b.channelManager.DeleteChannelInfo(chName)
		return errors.Errorf("Failed to create channel : The creator channel client creation"+
			"Failed : %v ", err)
	}

	//从通道中的peer节点中选一个，以其所在组织的身份创建通道
	_, creatorDomain := utils.SplitContainerName(peers[0])
	creatorOrg := getOrgFromNetManager(creatorDomain, b.netManager)
	if creatorOrg == nil {
		b.channelManager.DeleteChannelInfo(chName)
		return errors.New("Can't get a channel creator org ")
	}

	chClients[creatorDomain], err = ci.NewChannelClient("Admin", creatorOrg)
	defer chClients[creatorDomain].Close()
	if err != nil {
		b.channelManager.DeleteChannelInfo(chName)
		return errors.Errorf("Failed to create channel : The creator channel client creation"+
			"Failed : %v ", err)
	}

	//创建通道
	_, err = chClients[creatorDomain].CreateChannel()
	if err != nil {
		b.channelManager.DeleteChannelInfo(chName)
		return errors.Errorf("Failed to create channel : %v ", err)
	}

	//将peer节点加入通道
	for _, peer := range peers {
		_, domain := utils.SplitContainerName(peer)
		if client, ok := chClients[domain]; ok {

			err = client.JoinChannel(peer)
			if err != nil {
				return errors.Errorf("Failed to join channel : %v", err)
			}
		} else {
			opOrg := getOrgFromNetManager(domain, b.netManager)
			if opOrg == nil {
				return errors.New("Can't get org to creat channel client")
			}
			client, err = ci.NewChannelClient("Admin", opOrg)
			if err != nil {
				return errors.Errorf("Failed to creat client to join channel %v", err)
			}

			defer client.Close()
			chClients[domain] = client

			err = client.JoinChannel(peer)
			if err != nil {
				return errors.Errorf("Failed to join channel : %v", err)
			}
		}
	}

	//在peer节点加入后，刷新sdk config file
	err = channel.WriteSdkConfig(b.configurator, b.channelManager)
	if err != nil {
		return errors.Errorf("Failed to set up channel : Failed to write sdk config : %v", err)
	}

	//更新anchor节点，还有点问题
	for domain := range anchors {
		client, ok := chClients[domain]
		if ok {
			client.RefreshSdk()
		} else {
			opOrg := getOrgFromNetManager(domain, b.netManager)
			if opOrg == nil {
				return errors.New("Can't get org to creat channel client")
			}
			client, err = ci.NewChannelClient("Admin", opOrg)
			if err != nil {
				return errors.Errorf("Failed to creat client to update anchor:%v", err)
			}
			chClients[domain] = client
		}
		defer client.Close()
		_, err = client.UpdateAnchorPeers()
		if err != nil {

			return errors.Errorf("Failed to update anchor : %v", err)
		}
	}
	return nil
}

func (b *Blockchain) AddPeerToChannel(chName string, peers []string) error {
	var chClients = make(map[string]*channel.ChannelClient)

	//获取作为加入对象的channel info对象
	ci, ok := b.GetChannelManager().GetChannelInfo(chName)
	if !ok {
		return errors.New("Channel didn't exist!")
	}
	//创建建立channel client所需的sdkconfig
	err := channel.WriteSdkConfig(b.configurator, b.channelManager)
	if err != nil {
		return errors.Errorf("Failed to add peer: Failed to creat sdk config : %v", err)
	}

	for _, peer := range peers {
		_, domain := utils.SplitContainerName(peer)
		if client, ok := chClients[domain]; ok {
			err = client.JoinChannel(peer)
		} else {
			opOrg := getOrgFromNetManager(domain, b.netManager)
			if opOrg == nil {
				return errors.New("Can't get org to creat channel client")
			}
			client, err = ci.NewChannelClient("Admin", opOrg)
			if err != nil {
				return errors.Errorf("Failed to creat client to join channel:%v ", err)
			}

			chClients[domain] = client
			defer client.Close()

			err = client.JoinChannel(peer)
			if err != nil {
				return errors.Errorf("Failed to join channel: %v", err)
			}
		}
	}
	return nil
}

//Contract相关
func (b *Blockchain) CreatContract(config contract.ChaincodeConfig) (int, error) {
	cm := b.GetContractManager()
	ci, err := cm.NewChaincodeInfo(contract.ChaincodeConfig{
		ChannelName:      config.ChannelName,
		ChaincodeName:    config.ChaincodeName,
		ChaincodeDesc:    config.ChaincodeDesc,
		ChaincodeGoPath:  config.ChaincodeGoPath,
		ChaincodePath:    config.ChaincodePath,
		ChaincodeVersion: config.ChaincodeVersion,
		ChaincodeLang:    config.ChaincodeLang,
	})
	if err != nil {
		return -1, errors.Errorf("Failed to creat contract info : %v", err)
	}
	return ci.ChaincodeID(), nil
}

func (b *Blockchain) InstallContract(id int) ([]string, error) {
	ci, ok := b.GetContractManager().GetChaincodeByID(id)
	if !ok {
		return nil, errors.New("Chaincode didn't exist!")
	}

	//获取合约要安装的通道的信息
	ch, ok := b.channelManager.GetChannelInfo(ci.ChannelName())
	if !ok {
		return nil, errors.New("The channel that contract belong to didn't exist!")
	}

	err := channel.WriteSdkConfig(b.configurator, b.channelManager)
	if err != nil {
		return nil, errors.Errorf("Failed to install contract: Failed to creat sdk config : %v", err)
	}

	//在属于不同组织的节点上安装合约需要不同组织的cc client，因此需要创建多个client，用map来管理
	var orgMark = make(map[string]bool)
	//需要向多个节点安装链码，其中有些可能成功有些失败。这个切片上在哪些地址成功安装的返回值
	var ccInstalled []string
	//由于可能有部分节点安装失败，所以在安装过程中发生错误不会直接返回，而是将错误信息都添加到MultiError中，
	//继续下一个组织中的链码安装，并最后一并返回所有错误
	var multiError error = nil
	//要向目标通道中的每个组织安装链码
	for _, peer := range ch.GetPeers() {
		_, domain := utils.SplitContainerName(peer)

		//如果尝试过在某个组织上安装过链码，则会在map中创建新的一项。如果ok为false，则需要尝试
		//在该组织中安装链码，否则不用尝试
		_, ok := orgMark[domain]
		if !ok {
			//开始在某个组织中安装链码，对该组织进行标记
			orgMark[domain] = true
			//获取创建client所需的org对象
			opOrg := getOrgFromNetManager(domain, b.netManager)
			if opOrg == nil {
				multiError = errors.Errorf("%v ; Failed to get %v org to creat client",
					multiError, domain)
				continue
			}

			client, err := ci.NewChaincodeClient("Admin", opOrg)
			if err != nil {
				multiError = errors.Errorf("%v ; Failed to creat chaincode client for %v: %v",
					multiError, domain, err)
				continue
			}
			//client创建成功，开始安装链码
			resp, err := client.InstallChaincode()
			if err != nil {
				multiError = errors.Errorf("%v ; Filed to install chaincode for %v : %v",
					multiError, domain, err)
				continue
			}
			//将安装成功的response存储起来
			for _, tar := range resp {
				ccInstalled = append(ccInstalled, tar.Target)
			}

			client.Close()
		}
	}
	return ccInstalled, multiError
}

func (b *Blockchain) InstantiateCC(id int, args []string, upgrade bool) error {
	ci, ok := b.GetContractManager().GetChaincodeByID(id)
	if !ok {
		return errors.New("Chaincode didn't exist!")
	}

	//获取合约要安装的通道的信息
	ch, ok := b.channelManager.GetChannelInfo(ci.ChannelName())
	if !ok {
		return errors.New("The channel that contract belong to didn't exist!")
	}
	_, domain := utils.SplitContainerName(ch.GetPeers()[0])
	opOrg := getOrgFromNetManager(domain, b.netManager)
	if opOrg == nil {
		return errors.New("Failed to get org for chaincode client")
	}

	//创建client所需的sdkconfig
	err := channel.WriteSdkConfig(b.configurator, b.channelManager)
	if err != nil {
		return errors.Errorf("Failed to instantiate cc: Failed to creat sdk config : %v", err)
	}

	client, err := ci.NewChaincodeClient("Admin", opOrg)
	if err != nil {
		return errors.Errorf("Failed to creat cc client : %v ", err)
	}

	if upgrade {
		_, err = client.UpgradeChaincode(args...)
		if err != nil {
			return errors.Errorf("Failed to upgrade chaincode : %v", err)
		}
	} else {
		_, err = client.InstantiateChaincode(args...)
		if err != nil {
			return errors.Errorf("Failed to instantiate chaincode : %v", err)
		}
	}

	return nil
}

func (b *Blockchain) InvokeChaincode(id int, args []string, query bool) (CCInvokeResponseToJson, error) {
	ci, ok := b.contractManager.GetChaincodeByID(id)
	if !ok {
		return CCInvokeResponseToJson{}, errors.New("Chaincode didn't exist!")
	}

	//获取调用合约所在的通道
	ch, ok := b.channelManager.GetChannelInfo(ci.ChannelName())
	if !ok {
		return CCInvokeResponseToJson{}, errors.New("The channel that contract belong to didn't exist!")
	}

	_, domain := utils.SplitContainerName(ch.GetPeers()[0])
	opOrg := getOrgFromNetManager(domain, b.netManager)
	if opOrg == nil {
		return CCInvokeResponseToJson{}, errors.New("Can't get org to creat cliet")
	}

	err := channel.WriteSdkConfig(b.configurator, b.channelManager)
	if err != nil {
		return CCInvokeResponseToJson{}, errors.Errorf("Failed to add peer: Failed to creat sdk config : %v", err)
	}

	client, err := ci.NewChaincodeClient("Admin", opOrg)
	if err != nil {
		return CCInvokeResponseToJson{}, errors.Errorf("Failed to invoke chaincode : %v", err)
	}

	var result CCInvokeResponseToJson

	if query {
		resp, err := client.QueryChaincode(args[0], args[1:]...)
		b.ccLogger.Log(*ci, args, resp, query)
		if err != nil {
			return CCInvokeResponseToJson{}, errors.Errorf("Failed to invoke chaincode : %v", err)
		}
		result.Txid = string(resp.TransactionID)
		result.Status = resp.ChaincodeStatus
		result.Payload = string(resp.Payload)
	} else {
		resp, err := client.InvokeChaincode(args[0], args[1:]...)
		b.ccLogger.Log(*ci, args, resp, query)
		if err != nil {
			return CCInvokeResponseToJson{}, errors.Errorf("Failed to invoke chaincode : %v", err)
		}
		result.Txid = string(resp.TransactionID)
		result.Status = resp.ChaincodeStatus
		result.Payload = string(resp.Payload)
	}

	return result, nil
}

//区块链启动相关，应该由同模块内Blockchain Manager调用。
func (b *Blockchain) generateArtifacts() error {
	err := artifacts.GenerateYaml(b.configurator, utils.ConfigPathWithId(b.id))
	if err != nil {
		return errors.Errorf("Failed to prepare artifacts : %v ", err)
	}

	err = artifacts.GenerateCryptoConfig(utils.ConfigPathWithId(b.id))
	if err != nil {
		return errors.Errorf("Failed to prepare artifacts : %v ", err)
	}

	err = artifacts.GenerateGenesisBlock(utils.ConfigPathWithId(b.id))
	if err != nil {
		return errors.Errorf("Failed to prepare artifacts : %v ", err)
	}

	return nil
}

func (b *Blockchain) prepareBlockchain() error {
	//向netmanager加载网络结构
	err := b.netManager.LoadConfig(b.configurator)
	if err != nil {
		return errors.Errorf("Failed to Creat blockchain : Failed to load "+
			"configurator : %v", err)
	}

	fmt.Println("Prepared network")
	err = b.netManager.Prepared()
	if err != nil {
		return errors.Errorf("Failed to start blockchain net : %v ", err)
	}

	return nil
}

func (b *Blockchain) startBlockchain() error {

	if b.status == Configuring {
		return errors.New("blockchain not initialized")
	}

	if b.status == Running {
		return errors.New("blockchain is running")
	}

	fmt.Println("Startnet")
	err := b.netManager.StartNet()
	if err != nil {
		return errors.Errorf("Failed to start blockchain net : %v ", err)
	}
	b.status = Running

	return nil
}

func (b *Blockchain) stopBlockchain() error {
	if b.status != Running {
		return errors.New("blockchain isn't running")
	}

	err := b.netManager.StopNet()
	if err != nil {
		return errors.Errorf("Failed to stop blockchain net : %v ", err)
	}
	b.status = Stop

	return nil
}

func (b *Blockchain) removeBlockchain() error {
	if b.status == Running {
		return errors.New("Can't remove blockchain, " +
			"blockchain is running")
	}

	if b.status != Configuring {
		//调用network模块中的Remove方法，删除创建的docker镜像
		err := b.netManager.Remove()
		if err != nil {
			return errors.Errorf("Faild to RemoveBlockchain :%v ", err)
		}
	}

	//删除为Blockchain创建的文件
	err := os.RemoveAll(utils.ConfigPathWithId(b.id))
	if err != nil {
		return errors.Errorf("Failed to RemoveBlockchain :%v ", err)
	}
	return nil
}

//获取Log
func (b *Blockchain) GetContractLogs() []logger.ContractInvokeLog {
	return b.ccLogger.GetLog()
}

//从network manager中取出domain对应的组织，如果domain对应的组织并存在，则返回Nil
func getOrgFromNetManager(domain string, nm *network.NetManager) *network.PeerOrg {
	var result *network.PeerOrg
	result = nil
	for _, org := range nm.Peerorgs {
		if org.Domain() == domain {
			result = org
			break
		}
	}

	return result
}
