package channel

import (
	"bit-bass/network"
	"github.com/hyperledger/fabric-protos-go/common"
	"github.com/hyperledger/fabric-protos-go/peer"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/fab"
)

type BlockChainIf interface {
	GetBlockchainInfo() (*fab.BlockchainInfoResponse, error)
	GetBlockByHeight(height uint64) (*common.Block, error)
	GetBlcokByHash(hash []byte) (*common.Block, error)
	GetTransaction(txid string) (*peer.ProcessedTransaction, error)
}

type ChannelIf interface {
	BlockChainIf
	GetContracts(peer *network.PeerNode) (*peer.ChaincodeQueryResponse, error)
	ChannelName() string
	CreateChannel() (fab.TransactionID, error)
	JoinChannel(string) error
	UpdateAnchorPeers() (fab.TransactionID, error)
	RefreshSdk() error
	QueryChannelofPeer(*network.PeerNode) ([]string, error)
	Close()
}
