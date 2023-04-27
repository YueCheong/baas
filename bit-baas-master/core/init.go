package core

var CoreBlockChainManager *BlockchainManager

func InitCore() {
	CoreBlockChainManager = NewBlockchainManager()
}
