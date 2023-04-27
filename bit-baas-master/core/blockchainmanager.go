package core

import (
	"bit-bass/channel"
	"bit-bass/contract"
	"bit-bass/deploy"
	"bit-bass/logger"
	"bit-bass/network"
	"bit-bass/utils"
	"fmt"
	"github.com/docker/go-connections/nat"
	"github.com/pkg/errors"
	"os"
	"strconv"
)

/*
负责的功能：
管理区块链网络，分配ID
将任务委托给具体的区块链网络
*/
type BlockchainManager struct {
	blockchains map[int]*Blockchain
	dockernets  map[string]deploy.NetIf
	idGen       *utils.AutoIncIDGen
}

func NewBlockchainManager() *BlockchainManager {
	bm := BlockchainManager{}
	bm.blockchains = make(map[int]*Blockchain)
	bm.dockernets = make(map[string]deploy.NetIf)
	bm.idGen = utils.NewAutoIncID()

	return &bm
}

const NoBlockchainCreated = -1

//创建一个新的区块链
func (bm *BlockchainManager) NewBlockchain(config BlockchainToJson) (int, error) {
	id := bm.idGen.GenID()

	var configurator *deploy.Configurator
	var err error
	if config.Netname == "" {
		config.Netname = "fabric_net_" + strconv.Itoa(id)
		net, err := bm.NewDockerNet(config.Netname)
		if err != nil {
			return NoBlockchainCreated, errors.New("Can't creat net for blockchain")
		}
		configurator, err = genConfiguratorFromApiData(id, config, net)
	} else {
		net, ok := bm.GetNetByName(config.Netname)
		if !ok {
			return NoBlockchainCreated, errors.New("Can't creat blockchain, docker net didn't exist")
		}

		configurator, err = genConfiguratorFromApiData(id, config, net)
	}

	if err != nil {
		return NoBlockchainCreated, errors.Errorf("Failed to creat blockchain : Failed to creat configutar : %v ", err)
	}

	b := Blockchain{
		id:              id,
		name:            config.Name,
		netManager:      network.NewNetManager(),
		channelManager:  channel.NewChannelmanager(utils.ConfigPathWithId(id)),
		contractManager: contract.NewContractManager(utils.ConfigPathWithId(id)),
		configurator:    configurator,
		ccLogger:        logger.NewContractInvokeLogger(id, config.Name),
		status:          Configuring,
		dockerNetName:   config.Netname,
	}

	//将创建的区块链存储到manager中
	bm.blockchains[b.id] = &b

	return b.id, nil
}

func (bm *BlockchainManager) SetOrdererOrg(id int, config BlockchainManageInfoToJson) error {
	if b, ok := bm.blockchains[id]; ok {
		if b.configurator.Ordererorgconf != nil {
			return errors.New("Orderer org already exist.")
		}

		if config.Name != "" && config.Domain != "" && config.MSPID != "" {
			err := b.configurator.SetOrdererOrg(config.Name, config.Domain, config.MSPID)
			if err != nil {
				return errors.Errorf("Failed to set configurator : %v", err)
			}
		} else {
			return errors.New("Org name, domain or MSPID can't be empty")
		}

		return nil
	}
	return errors.New("Can't find blockchain with specified id")
}

func (bm *BlockchainManager) SetOrderer(id int, config BlockchainManageInfoToJson) error {
	if b, ok := bm.blockchains[id]; ok {
		if b.configurator.Ordererorgconf == nil {
			return errors.New("There's no orderer org")
		}

		if b.configurator.Ordererorgconf.GetNodes()[0] != nil {
			return errors.New("Orderer node already exist")
		}

		fmt.Println(config.Nodes)

		if len(config.Nodes) != 0 {
			node := config.Nodes[0]
			err := b.configurator.SetOrdererNode(node.Name, nat.Port(node.Port))
			if err != nil {
				return errors.Errorf("Failed to set configurator : %v", err)
			}
		} else {
			return errors.New("There's no orderer node")
		}

		return nil
	}
	return errors.New("Can't find blockchain with specified id")
}

func (bm *BlockchainManager) AddPeerOrg(id int, config BlockchainManageInfoToJson) error {
	if b, ok := bm.blockchains[id]; ok {

		if config.Name != "" && config.Domain != "" && config.MSPID != "" {
			err := b.configurator.AddPeerOrg(config.Name, config.Name+"."+config.Domain, config.MSPID)
			if err != nil {
				return errors.Errorf("Failed to set configurator : %v", err)
			}
		} else {
			return errors.New("Org name, domain or MSPID can't be empty")
		}

		return nil
	}
	return errors.New("Can't find blockchain with specified id")
}

func (bm *BlockchainManager) AddPeer(id int, config BlockchainManageInfoToJson) error {
	if b, ok := bm.blockchains[id]; ok {

		for _, peer := range config.Nodes {
			err := b.configurator.AddPeerNodeToOrg(config.Name, peer.Name, nat.Port(peer.Port))
			if err != nil {
				return errors.Errorf("Failed to set configurator : %v", err)
			}
		}

		return nil
	}
	return errors.New("Can't find blockchain with specified id")
}

func (bm *BlockchainManager) InitializeBlockchainById(id int) error {
	if b, ok := bm.blockchains[id]; ok {
		if b.status != Configuring {
			return errors.New("Blockchain is already initialized")
		}

		//检查系统是否有最基础的网络架构
		//有Orderer org和orderer
		if b.configurator.Ordererorgconf == nil {
			return errors.New("Can't initialize Blockchain: Lack of orderer org")
		}

		if b.configurator.Ordererorgconf.GetNodes()[0] == nil {
			return errors.New("Can't initialize Blockchain: Lack of orderer node")
		}
		//至少有一个peer org和peer
		if len(b.configurator.Peerorgsconf) == 0 {
			return errors.New("Can't initialize Blockchain: Lack of peer org")
		}

		for _, peerOrgConf := range b.configurator.Peerorgsconf {
			if len(peerOrgConf.GetNodes()) == 0 {
				return errors.New("Can't initialize Blockchain: There's peer org without peer node")
			}
		}

		//添加client
		err := b.configurator.SetToolsCli("cli", b.configurator.Peerorgsconf[0].GetNodes()[0])
		if err != nil {
			return errors.Errorf("Failed to set cli config : %v", err)
		}

		//创建所需的文件
		err = b.generateArtifacts()
		if err != nil {
			errf := os.RemoveAll(utils.ConfigPathWithId(b.id))
			if errf != nil {
				fmt.Println("区块链创建失败，且产生的文件无法被删除")
			}
			return errors.Errorf("Failed to creat blockchain : Failed to creat "+
				"artifacts : %v ", err)
		}

		//区块链启动前的准备
		err = b.prepareBlockchain()
		if err != nil {
			err = os.RemoveAll(utils.ConfigPathWithId(b.id))
			if err != nil {
				fmt.Println("区块链创建失败，且产生的文件无法被删除")
			}
			return errors.Errorf("Failed to creat blockchain : Failed to creat"+
				"artifacts : %v ", err)
		}

		b.status = Stop

		return nil
	}

	return errors.New("Blockchain not found.")
}

//创建一个网络
func (bm *BlockchainManager) NewDockerNet(name string) (*network.NodesNet, error) {
	_, ok := bm.dockernets[name]
	if ok {
		return nil, errors.New("Network with same name already exist")
	}

	net := network.NewDefaultNodesNet(name)
	if net == nil {
		return nil, errors.New("Failed to creat docker net item")
	}

	err := net.CreateNet()
	if err != nil {
		return nil, errors.Errorf("Failed to creat docker net : %v ", err)
	}

	bm.dockernets[name] = net
	return net, nil
}

//根据网络名获取网络
func (bm *BlockchainManager) GetNetByName(name string) (deploy.NetIf, bool) {
	netI, ok := bm.dockernets[name]
	return netI, ok
}

//获取所有网络
func (bm *BlockchainManager) GetNets() []deploy.NetIf {
	var result []deploy.NetIf
	for _, value := range bm.dockernets {
		result = append(result, value)
	}

	return result
}

//根据网络名删除网络
func (bm *BlockchainManager) RemoveNetByName(name string) error {
	if net, ok := bm.dockernets[name]; ok {
		err := net.RemoveNet()
		if err != nil {
			return errors.Errorf("Failed to remove docker net : %v ", err)
		}

		delete(bm.dockernets, name)
		return nil
	}

	return errors.New("Net didn't exist")
}

//依照ID启动区块链系统
func (bm *BlockchainManager) StartBlockchainById(id int) error {
	if b, ok := bm.blockchains[id]; ok {
		err := b.startBlockchain()
		if err != nil {
			return errors.Errorf("Failed to start blockchain : %v ", err)
		}
		return nil
	}

	return errors.New("Can't find blockchain with specified id")

}

//依照id停止区块链系统
func (bm *BlockchainManager) StopBlockchainById(id int) error {
	if b, ok := bm.blockchains[id]; ok {
		err := b.stopBlockchain()
		if err != nil {
			return errors.Errorf("Failed to stop blockchain : %v", err)
		}
		return nil
	}
	return errors.New("Can't find blockchain with specified id")
}

//依照id删除区块链系统
func (bm *BlockchainManager) RemoveBlockchainById(id int) error {
	if b, ok := bm.blockchains[id]; ok {
		err := b.removeBlockchain()
		if err != nil {
			return errors.Errorf("Failed to remove blockchain : %v ", err)
		}

		delete(bm.blockchains, id)

		return nil
	}
	return errors.New("Can't find blockchain with specified id")
}

//获取所有的区块链
func (bm *BlockchainManager) GetBlockchains() []*Blockchain {
	var result []*Blockchain
	for k := range bm.blockchains {
		result = append(result, bm.blockchains[k])
	}

	return result
}

//依照ID获取区块链
func (bm *BlockchainManager) GetBlockchainById(id int) (*Blockchain, bool) {
	b, ok := bm.blockchains[id]
	return b, ok
}

//依照ID关闭并清除区块链
func (bm *BlockchainManager) StopAndRemoveBlockchainById(id int) error {
	if b, ok := bm.blockchains[id]; ok {
		err := b.stopBlockchain()
		if err != nil {
		}
		err = b.removeBlockchain()
		if err != nil {
			return errors.Errorf("Failed to delete blockchain : %v", err)
		}

		delete(bm.blockchains, id)

		return nil
	}

	return errors.New("Can't find blockchain with specified id")
}

func genConfiguratorFromApiData(id int, config BlockchainToJson, net deploy.NetIf) (*deploy.Configurator, error) {
	c := deploy.NewConfigurator(utils.ConfigPathWithId(id), net)

	if config.OrdererOrg.Name != "" && config.OrdererOrg.Domain != "" &&
		config.OrdererOrg.MSPID != "" {
		err := c.SetOrdererOrg(config.OrdererOrg.Name, config.OrdererOrg.Domain, config.OrdererOrg.MSPID)
		if err != nil {
			return nil, errors.Errorf("Failed to set configurator : %v", err)
		}

		if config.OrdererOrg.Orderer.Name != "" && config.OrdererOrg.Orderer.Port != "" {
			err = c.SetOrdererNode(config.OrdererOrg.Orderer.Name, nat.Port(config.OrdererOrg.Orderer.Port))
			if err != nil {
				return nil, errors.Errorf("Failed to set configurator : %v", err)
			}
		}
	}

	for _, peerOrg := range config.PeerOrg {
		err := c.AddPeerOrg(peerOrg.Name, peerOrg.Name+"."+peerOrg.Domain, peerOrg.MSPID)
		if err != nil {
			return nil, errors.Errorf("Failed to set configurator : %v", err)
		}
		//
		//// 每个组织都应该只有一个anchor peer
		//err = c.SetAnchorPeerToOrg(peerOrg.Name, peerOrg.AnchorPeer)
		//if err != nil {
		//	return nil, errors.Errorf("Failed to set configurator : %v", err)
		//}

		for _, peer := range peerOrg.Peers {
			err = c.AddPeerNodeToOrg(peerOrg.Name, peer.Name, nat.Port(peer.Port))
			if err != nil {
				return nil, errors.Errorf("Failed to set configurator : %v", err)
			}
		}
	}

	return c, nil
}

func (bm *BlockchainManager) GenerateSummary() (*BaasSummaryToJson, error) {
	blockchains := CoreBlockChainManager.GetBlockchains()
	var blockTotal uint64
	result := BaasSummaryToJson{
		BlochainTotal:     len(blockchains),
		BlockchainRunning: 0,
	}

	for _, b := range blockchains {
		bSum := BlockchainSummary{
			BlockchainName: b.GetName(),
			BlockchainID:   b.GetId(),
			OrdererOrgNum:  1,
			PeerOrgNum:     len(b.GetConfigurator().Peerorgsconf),
			ChannelNum:     len(b.GetChannelManager().GetChannelinfos()),
			ContractNum:    len(b.GetContractManager().GetChaincodes()),
		}
		if b.GetStatus() == Running {
			result.BlockchainRunning++
		}

		//fill b.ChannelSummarys
		for _, ch := range b.GetChannelManager().GetChannelinfos() {
			chSum := ChannelSummary{
				ChannelName: ch.GetChannelName(),
			}
			_, domain := utils.SplitContainerName(ch.GetPeers()[0])
			chCli, err := ch.NewChannelClient("Admin", getOrgFromNetManager(domain, b.netManager))
			if err != nil {
				return nil, errors.New("Can't creat ch client")
			}
			defer chCli.Close()

			infoResp, err := chCli.GetBlockchainInfo()
			if err != nil {
				chSum.Blockheight = "Failed to get block height :" + err.Error()
			} else {
				chSum.Blockheight = strconv.FormatUint(infoResp.BCI.Height, 10)
				blockTotal += infoResp.BCI.Height
			}

			bSum.ChannelSummarys = append(bSum.ChannelSummarys, chSum)
		}
		result.BlockchainSummarys = append(result.BlockchainSummarys, bSum)
	}

	for _, bSum := range result.BlockchainSummarys {
		result.OrganizationTotal += (bSum.OrdererOrgNum + bSum.PeerOrgNum)
		result.OrdererOrgTotal += bSum.OrdererOrgNum
		result.PeerOrgTotal += bSum.PeerOrgNum
		result.ChannelTotal += bSum.ChannelNum
		result.ContractTotal += bSum.ContractNum
		result.BlockTotal = strconv.FormatUint(blockTotal, 10)
	}

	return &result, nil
}
