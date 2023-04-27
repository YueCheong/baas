package channel

import (
	"bit-bass/artifacts"
	"bit-bass/network"
	"github.com/hyperledger/fabric-protos-go/common"
	"github.com/hyperledger/fabric-protos-go/peer"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/ledger"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/resmgmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/fab"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
	"github.com/pkg/errors"
	"os"
)

type ChannelClient struct {
	info        *Channelinfo
	operateOrg  network.OrgIf
	operateUser string
	sdk         *fabsdk.FabricSDK
}

func (c *ChannelClient) GetBlockchainInfo() (*fab.BlockchainInfoResponse, error) {
	ledgerCli, err := c.getLedgerClient()
	if err != nil {
		return nil, errors.Errorf("Failed to get Ledger Client : %v", err)
	}

	blockchainInfo, err := ledgerCli.QueryInfo()
	if err != nil {
		return nil, errors.Errorf("Failed to get blockchain info : %v", err)
	}

	return blockchainInfo, nil
}

func (c *ChannelClient) GetBlockByHeight(height uint64) (*common.Block, error) {
	//get ledger client
	ledgerCli, err := c.getLedgerClient()
	if err != nil {
		return nil, errors.Errorf("Failed to get Ledger Client : %v", err)
	}

	block, err := ledgerCli.QueryBlock(height)
	if err != nil {
		return nil, errors.Errorf("Failed to get block : %v", err)
	}

	return block, nil
}

func (c *ChannelClient) GetBlcokByHash(hash []byte) (*common.Block, error) {
	ledgerCli, err := c.getLedgerClient()
	if err != nil {
		return nil, errors.Errorf("Failed to get Ledger Client : %v", err)
	}

	block, err := ledgerCli.QueryBlockByHash(hash)
	if err != nil {
		return nil, errors.Errorf("Failed to get block : %v", err)
	}

	return block, nil
}

func (c *ChannelClient) GetTransaction(txid string) (*peer.ProcessedTransaction, error) {
	ledgerCli, err := c.getLedgerClient()
	if err != nil {
		return nil, errors.Errorf("Failed to get Ledger Client : %v", err)
	}

	tx, err := ledgerCli.QueryTransaction(fab.TransactionID(txid))
	if err != nil {
		return nil, errors.Errorf("Failed to get transaction info : %v", err)
	}

	return tx, nil
}

func (c *ChannelClient) GetContracts(peer *network.PeerNode) (*peer.ChaincodeQueryResponse, error) {
	//Creat res Client
	resClient, err := resmgmt.New(c.sdk.Context(fabsdk.WithOrg(c.operateOrg.Name()), fabsdk.WithUser(c.operateUser)))
	if err != nil {
		return nil, errors.Errorf("Failed to creat resClient : %v", err)
	}

	resp, err := resClient.QueryInstalledChaincodes(resmgmt.WithTargetEndpoints(peer.ContainerName()))
	if err != nil {
		return nil, errors.Errorf("Failed to query chaincodes : %v", err)
	}

	return resp, nil
}

func (c *ChannelClient) ChannelName() string {
	return c.info.channelName
}

func (c *ChannelClient) CreateChannel() (fab.TransactionID, error) {
	chConfigPath, err := artifacts.GenerateChannelCreationTx(c.info.configPath, c.info.channelName)
	if err != nil {
		return "", errors.Errorf("Failed to create channel : Failed to creat tx : %v", err)
	}
	//Creat channel creat request
	channelConfig, err := os.Open(chConfigPath)
	if err != nil {
		return "", errors.Errorf("Failed to open channel artifacts: %s", err)
	}
	defer channelConfig.Close()

	req := resmgmt.SaveChannelRequest{
		ChannelID:     c.info.channelName,
		ChannelConfig: channelConfig,
	}

	//Propose ChannelClient creat request
	response, err := c.proposeChannelChangeTransaction(req)

	if err != nil {
		return "", errors.Errorf("Error from save channel :%s", err)
	}
	if response.TransactionID == "" {
		return "", errors.Errorf("Failed to save channel")
	}

	c.info.created = true

	return response.TransactionID, nil
}

func (c *ChannelClient) JoinChannel(peerContainerName string) error {

	//Creat res Client with org and user
	resClient, err := c.getResmgmtClient()
	if err != nil {
		return err
	}

	err = resClient.JoinChannel(c.info.channelName,
		resmgmt.WithTargetEndpoints(peerContainerName))
	if err != nil {
		return errors.Errorf("Failed to join channel in ChannelClient.JoinChannel :%v", err)
	}
	c.info.peers = append(c.info.peers, peerContainerName)
	return nil
}

func (c *ChannelClient) UpdateAnchorPeers() (fab.TransactionID, error) {

	txpath, err := artifacts.GenerateAnchorPeerTx(c.info.configPath, c.info.channelName, c.operateOrg.Name())
	if err != nil {
		return "", errors.Errorf("Failed to update anchor peer : Failed to creat tx : %v ", err)
	}
	//Creat channel AnchorPeer Update request
	channelConfig, err := os.Open(txpath)
	if err != nil {
		return "", errors.Errorf("Failed to open channel artifacts: %s", err)
	}
	defer channelConfig.Close()

	req := resmgmt.SaveChannelRequest{
		ChannelID:     c.info.channelName,
		ChannelConfig: channelConfig,
	}

	//Propose ChannelClient creat request
	response, err := c.proposeChannelChangeTransaction(req)

	if err != nil {
		return "", errors.Errorf("Error from save channel :%s", err)
	}
	if response.TransactionID == "" {
		return "", errors.Errorf("Failed to save channel")
	}

	return response.TransactionID, nil
}

//当区块连网络发生变化时，如节点改变或创建新通道，新节点加入通道等，需要重新生成新的sdk配置文件，此函数的作用是根据新到
//配置文件更新Channel对象中到sdk
func (c *ChannelClient) RefreshSdk() error {
	c.sdk.Close()
	netConifg := config.FromFile(c.info.configPath + "/sdk-config.yaml")
	fabricsdk, err := fabsdk.New(netConifg)
	if err != nil {
		return errors.Errorf("Faild to creat fabric sdk : %v", err)
	}

	c.sdk = fabricsdk
	return nil
}

func (c *ChannelClient) QueryChannelofPeer(peer *network.PeerNode) ([]string, error) {
	resCli, err := c.getResmgmtClient()
	if err != nil {
		return nil, errors.Errorf("Failed to creat Resmgmt Client : %v", err)
	}
	resp, err := resCli.QueryChannels(resmgmt.WithTargetEndpoints(peer.ContainerName()))
	if err != nil {
		return nil, errors.Errorf("Failed to query channels : %v", err)
	}

	var channelIDs []string
	for _, channel := range resp.GetChannels() {
		channelIDs = append(channelIDs, channel.GetChannelId())
	}

	return channelIDs, nil
}

func (c *ChannelClient) Close() {
	c.sdk.Close()
}

//把Channel状态变更请求发布到区块连中
func (c *ChannelClient) proposeChannelChangeTransaction(req resmgmt.SaveChannelRequest) (resmgmt.SaveChannelResponse, error) {
	emptyResp := resmgmt.SaveChannelResponse{}

	//Creat res Client
	resClient, err := c.getResmgmtClient()
	if err != nil {
		return emptyResp, errors.Errorf("Faild to creat res client :%v", err)
	}

	//propose the transaction
	resp, err := resClient.SaveChannel(req)
	if err != nil {
		return resp, err
	}

	return resp, nil
}

func (c *ChannelClient) getLedgerClient() (*ledger.Client, error) {
	ledgerCli, err := ledger.New(c.sdk.ChannelContext(c.info.channelName,
		fabsdk.WithUser(c.operateUser), fabsdk.WithOrg(c.operateOrg.Name())))

	if err != nil {
		return ledgerCli, errors.Errorf("Failed to creat ledger client : %v", err)
	}

	return ledgerCli, nil
}

func (c *ChannelClient) getResmgmtClient() (*resmgmt.Client, error) {

	resClient, err := resmgmt.New(c.sdk.Context(fabsdk.WithUser(c.operateUser), fabsdk.WithOrg(c.operateOrg.Name())))
	if err != nil {
		return nil, errors.Errorf("Failed to creat resmgmt client :%v", err)
	}

	return resClient, nil
}
