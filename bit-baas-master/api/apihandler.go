package api

import (
	"bit-bass/contract"
	"bit-bass/core"
	"bit-bass/logger"
	"bit-bass/utils"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

//-----------------------
//此处开始为handler函数

func GetNetworks(c *gin.Context) {
	netIs := core.CoreBlockChainManager.GetNets()
	var result []string
	for _, net := range netIs {
		result = append(result, net.NetName())
	}

	c.JSON(SUCCESS.Int(), Response{
		Code:    SUCCESS,
		Msg:     SUCCESS.GetMsg(),
		Package: result,
	})
}

func CreatNetwork(c *gin.Context) {
	name := c.Param("name")
	if name == "" {
		c.JSON(INVALID_PARAMS.Int(), Response{
			Code: INVALID_PARAMS,
			Msg:  INVALID_PARAMS.GetMsg(),
		})
		return
	}

	_, err := core.CoreBlockChainManager.NewDockerNet(name)
	if err != nil {
		c.JSON(ERROR.Int(), Response{
			Code:    DOCKER_NETWORK_ERROR,
			Msg:     DOCKER_NETWORK_ERROR.GetMsg(),
			Package: err.Error(),
		})
		return
	}

	c.JSON(SUCCESS.Int(), gin.H{})
}

func RemoveNetwork(c *gin.Context) {
	name := c.Param("name")
	if name == "" {
		c.JSON(INVALID_PARAMS.Int(), Response{
			Code: INVALID_PARAMS,
			Msg:  INVALID_PARAMS.GetMsg(),
		})
		return
	}

	err := core.CoreBlockChainManager.RemoveNetByName(name)
	if err != nil {
		c.JSON(ERROR.Int(), Response{
			Code:    DOCKER_NETWORK_ERROR,
			Msg:     DOCKER_NETWORK_ERROR.GetMsg(),
			Package: err.Error(),
		})
		return
	}

	c.JSON(SUCCESS.Int(), gin.H{})
}

//获取整个系统的概览信息
func GetSummary(c *gin.Context) {
	summary, err := core.CoreBlockChainManager.GenerateSummary()
	if err != nil {
		c.JSON(ERROR.Int(), Response{
			Code:    ERROR,
			Msg:     ERROR.GetMsg(),
			Package: err.Error(),
		})
		return
	}

	c.JSON(SUCCESS.Int(), Response{
		Code:    SUCCESS,
		Msg:     SUCCESS.GetMsg(),
		Package: summary,
	})
}

//获取系统中所有的区块链
func GetBlockchains(c *gin.Context) {
	blockchains := core.CoreBlockChainManager.GetBlockchains()
	var result []core.BlockchainToJson

	for _, b := range blockchains {
		result = append(result, getBlockchainToJsonFromBlockchain(b))
	}

	c.JSON(SUCCESS.Int(), Response{
		Code:    SUCCESS,
		Msg:     SUCCESS.GetMsg(),
		Package: result,
	})
}

//创建一个新的区块链网络，并创建请求中附带的通道
func CreatBlockchain(c *gin.Context) {
	var config core.BlockchainToJson
	err := c.BindJSON(&config)

	if err != nil {
		c.JSON(INVALID_PARAMS.Int(), Response{
			Code:    INVALID_PARAMS,
			Msg:     INVALID_PARAMS.GetMsg(),
			Package: err.Error(),
		})
		return
	}

	_, err = core.CoreBlockChainManager.NewBlockchain(config)
	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			Code:    BLOCKCHAIN_NETWORK_ERROR,
			Msg:     BLOCKCHAIN_NETWORK_ERROR.GetMsg(),
			Package: err.Error(),
		})
		return
	}

	c.JSON(SUCCESS.Int(), Response{
		Code: SUCCESS,
		Msg:  "区块链已被成功创建并初始化",
	})
}

//管理网络，包括网络的停止和启动
/*
postform 参数
Operation   int start = 1 stop = 2
*/

func ManageBlockchain(c *gin.Context) {
	var param core.BlockchainManageInfoToJson
	err := c.BindJSON(&param)
	if err != nil {
		c.JSON(INVALID_PARAMS.Int(), Response{
			Code:    INVALID_PARAMS,
			Msg:     INVALID_PARAMS.GetMsg(),
			Package: err.Error(),
		})
		return
	}

	_, ok := core.CoreBlockChainManager.GetBlockchainById(param.BlockchainID)
	if !ok {
		c.JSON(http.StatusInternalServerError, Response{
			Code:    RESOURCES_NOT_FOUND,
			Msg:     RESOURCES_NOT_FOUND.GetMsg(),
			Package: "找不到请求的id对应的区块连网络",
		})
		return
	}

	switch param.Operation {
	case core.StartBlockchain:
		err = core.CoreBlockChainManager.StartBlockchainById(param.BlockchainID)
	case core.StopBlockchain:
		err = core.CoreBlockChainManager.StopBlockchainById(param.BlockchainID)
	case core.InitializeBlockchain:
		err = core.CoreBlockChainManager.InitializeBlockchainById(param.BlockchainID)
	case core.SetOrdererOrg:
		err = core.CoreBlockChainManager.SetOrdererOrg(param.BlockchainID, param)
	case core.SetOrderer:
		err = core.CoreBlockChainManager.SetOrderer(param.BlockchainID, param)
	case core.AddPeerOrg:
		err = core.CoreBlockChainManager.AddPeerOrg(param.BlockchainID, param)
	case core.AddPeer:
		err = core.CoreBlockChainManager.AddPeer(param.BlockchainID, param)
	default:
		c.JSON(INVALID_PARAMS.Int(), Response{
			Code:    INVALID_PARAMS,
			Msg:     "Unknown Operation",
			Package: nil,
		})
		return
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			Code:    BLOCKCHAIN_NETWORK_ERROR,
			Msg:     BLOCKCHAIN_NETWORK_ERROR.GetMsg(),
			Package: err.Error(),
		})
		return
	}

	c.JSON(SUCCESS.Int(), Response{
		Code:    SUCCESS,
		Msg:     SUCCESS.GetMsg(),
		Package: nil,
	})
}

//关闭并删除区块链网络
func DeleteBlockchain(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(INVALID_PARAMS.Int(), Response{
			Code:    INVALID_PARAMS,
			Msg:     INVALID_PARAMS.GetMsg(),
			Package: err.Error(),
		})
		return
	}

	err = core.CoreBlockChainManager.StopBlockchainById(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			Code:    BLOCKCHAIN_NETWORK_ERROR,
			Msg:     BLOCKCHAIN_NETWORK_ERROR.GetMsg(),
			Package: err.Error(),
		})
		return
	}

	err = core.CoreBlockChainManager.RemoveBlockchainById(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			Code:    BLOCKCHAIN_NETWORK_ERROR,
			Msg:     BLOCKCHAIN_NETWORK_ERROR.GetMsg(),
			Package: err.Error(),
		})
		return
	}

	c.JSON(SUCCESS.Int(), gin.H{})
}

//获取某个区块链网络中所有通道的信息
func GetChannels(c *gin.Context) {
	idParam := c.Query("blockchainid")
	var blockchains []*core.Blockchain
	if idParam == "" {
		blockchains = core.CoreBlockChainManager.GetBlockchains()
	} else {
		id, err := strconv.Atoi(idParam)
		if err != nil {
			c.JSON(INVALID_PARAMS.Int(), Response{
				Code:    INVALID_PARAMS,
				Msg:     INVALID_PARAMS.GetMsg(),
				Package: err.Error(),
			})
			return
		}

		b, ok := core.CoreBlockChainManager.GetBlockchainById(id)
		if !ok {
			c.JSON(http.StatusInternalServerError, Response{
				Code:    RESOURCES_NOT_FOUND,
				Msg:     RESOURCES_NOT_FOUND.GetMsg(),
				Package: "找不到请求的id对应的区块连网络",
			})
			return
		}
		blockchains = append(blockchains, b)
	}

	var result []core.ChannelToJson
	for _, b := range blockchains {
		result = append(result, getChannelToJsonsFromBlockchain(b)...)
	}

	c.JSON(SUCCESS.Int(), Response{
		Code:    SUCCESS,
		Msg:     SUCCESS.GetMsg(),
		Package: result,
	})
}

func CreatChannel(c *gin.Context) {
	var ch core.ChannelToJson
	err := c.BindJSON(&ch)

	if err != nil {
		c.JSON(INVALID_PARAMS.Int(), Response{
			Code:    INVALID_PARAMS,
			Msg:     INVALID_PARAMS.GetMsg(),
			Package: err.Error(),
		})
		return
	}

	b, ok := core.CoreBlockChainManager.GetBlockchainById(ch.BlockchainID)
	if !ok {
		c.JSON(http.StatusInternalServerError, Response{
			Code:    RESOURCES_NOT_FOUND,
			Msg:     RESOURCES_NOT_FOUND.GetMsg(),
			Package: "找不到请求的id对应的区块连网络",
		})
		return
	}

	err = b.CreatAndInitChannel(ch.Name, ch.Peers, ch.AnchorPeers)
	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			Code:    BLOCKCHAIN_CHANNEL_ERROR,
			Msg:     BLOCKCHAIN_CHANNEL_ERROR.GetMsg(),
			Package: err.Error(),
		})
		return
	}

	c.JSON(SUCCESS.Int(), Response{
		Code:    SUCCESS,
		Msg:     SUCCESS.GetMsg(),
		Package: nil,
	})
}

func ManageChannel(c *gin.Context) {
	var config core.ChannelManageInfoToJson
	err := c.BindJSON(&config)
	if err != nil {
		c.JSON(INVALID_PARAMS.Int(), Response{
			Code:    INVALID_PARAMS,
			Msg:     INVALID_PARAMS.GetMsg(),
			Package: err.Error(),
		})
		return
	}

	b, ok := core.CoreBlockChainManager.GetBlockchainById(config.BlockchainID)
	if !ok {
		c.JSON(http.StatusInternalServerError, Response{
			Code:    RESOURCES_NOT_FOUND,
			Msg:     RESOURCES_NOT_FOUND.GetMsg(),
			Package: "找不到请求的id对应的区块连网络",
		})
		return
	}

	switch config.Operation {
	case core.AddPeerToChannel:
		addPeers(config.ChannelName, b, config.Args, c)
	case core.SetAnchorPeer:
		c.JSON(NOT_IMPLEMENTED.Int(), Response{
			Code:    NOT_IMPLEMENTED,
			Msg:     NOT_IMPLEMENTED.GetMsg(),
			Package: nil,
		})
	default:
		c.JSON(INVALID_PARAMS.Int(), Response{
			Code:    INVALID_PARAMS,
			Msg:     INVALID_PARAMS.GetMsg(),
			Package: "未知操作类型",
		})
	}
}

func GetContracts(c *gin.Context) {
	idParam := c.Query("blockchainid")
	chname := c.Query("channelname")

	var blockchains []*core.Blockchain
	var result []core.ContractToJson

	if idParam == "" {
		blockchains = core.CoreBlockChainManager.GetBlockchains()
	} else {
		id, err := strconv.Atoi(idParam)
		if err != nil {
			c.JSON(INVALID_PARAMS.Int(), Response{
				Code:    INVALID_PARAMS,
				Msg:     INVALID_PARAMS.GetMsg(),
				Package: err.Error(),
			})
			return
		}

		b, ok := core.CoreBlockChainManager.GetBlockchainById(id)
		if !ok {
			c.JSON(http.StatusInternalServerError, Response{
				Code:    RESOURCES_NOT_FOUND,
				Msg:     RESOURCES_NOT_FOUND.GetMsg(),
				Package: "找不到请求的id对应的区块连网络",
			})
			return
		}
		blockchains = append(blockchains, b)
	}

	if idParam != "" && chname != "" {
		for _, b := range blockchains {
			contracts := getContractToJsonsFromBlockchain(b)
			for _, cc := range contracts {
				if cc.ChannelName == chname {
					result = append(result, cc)
				}
			}
		}
	} else {
		for _, b := range blockchains {
			result = append(result, getContractToJsonsFromBlockchain(b)...)
		}
	}

	c.JSON(SUCCESS.Int(), Response{
		Code:    SUCCESS,
		Msg:     SUCCESS.GetMsg(),
		Package: result,
	})

}

/*
	postform 参数：
	BlockchainID   int
	ContractLang	int 0 = Golang 1 = java 2 = node
	ChannelName		string
	ContractName 	string
	ContractVersion	string
	ContractDesc	string
	file			cc source file
*/
func CreatContract(c *gin.Context) {

	//获得并校验输入

	bid, err := strconv.Atoi(c.PostForm("BlockchainID"))
	if err != nil {
		c.JSON(INVALID_PARAMS.Int(), Response{
			Code:    INVALID_PARAMS,
			Msg:     INVALID_PARAMS.GetMsg(),
			Package: err.Error(),
		})
		return
	}

	cclangInt, err := strconv.Atoi(c.PostForm("ContractLang"))
	cclang := contract.ChaincodeLanguageType(cclangInt)
	if err != nil || !cclang.Valid() {
		if err == nil {
			err = errors.New("Unknown chaincode language type")
		}
		c.JSON(INVALID_PARAMS.Int(), Response{
			Code:    INVALID_PARAMS,
			Msg:     INVALID_PARAMS.GetMsg(),
			Package: err.Error(),
		})
		return
	}

	chName := c.PostForm("ChannelName")
	if chName == "" {
		c.JSON(INVALID_PARAMS.Int(), Response{
			Code:    INVALID_PARAMS,
			Msg:     INVALID_PARAMS.GetMsg(),
			Package: "error : Channel name is empty",
		})
		return
	}

	ccName := c.PostForm("ContractName")
	if ccName == "" {
		c.JSON(INVALID_PARAMS.Int(), Response{
			Code:    INVALID_PARAMS,
			Msg:     INVALID_PARAMS.GetMsg(),
			Package: "error : Chaincode name is empty",
		})
		return
	}

	ccVer := c.PostForm("ContractVersion")
	if ccVer == "" {
		c.JSON(INVALID_PARAMS.Int(), Response{
			Code:    INVALID_PARAMS,
			Msg:     INVALID_PARAMS.GetMsg(),
			Package: "error : Chaincode version is empty",
		})
		return
	}

	ccDesc := c.PostForm("ContractDesc")

	b, ok := core.CoreBlockChainManager.GetBlockchainById(bid)
	if !ok {
		c.JSON(http.StatusInternalServerError, Response{
			Code:    RESOURCES_NOT_FOUND,
			Msg:     RESOURCES_NOT_FOUND.GetMsg(),
			Package: "找不到请求的id对应的区块连网络",
		})
		return
	}

	//上传并存储合约源代码
	//获得上传的文件
	form, err := c.MultipartForm()
	if err != nil || form == nil {
		c.JSON(ERROR.Int(), Response{
			Code:    ERROR,
			Msg:     ERROR.GetMsg(),
			Package: "Failed to get request form",
		})
		return
	}
	files := form.File["file"]
	if len(files) == 0 {
		c.JSON(INVALID_PARAMS.Int(), Response{
			Code:    INVALID_PARAMS,
			Msg:     INVALID_PARAMS.GetMsg(),
			Package: "No source file uploaded",
		})
		return
	}

	//生成链码存储文件名
	dst := utils.ConfigPathWithId(bid) + "/chaincode/"
	var ccpath string

	switch cclang {
	case contract.Golang:
		ccpath = ccName + "_" + ccVer + "_" + strconv.FormatInt(time.Now().Unix(), 10)
		dst = dst + "go/src/" + ccpath

	case contract.Java:
		dst = dst + "java/" + ccName + "_" + ccVer + "_" + strconv.FormatInt(time.Now().Unix(), 10)
		ccpath = dst
	case contract.Node:
		dst = dst + "node/" + ccName + "_" + ccVer + "_" + strconv.FormatInt(time.Now().Unix(), 10)
		ccpath = dst
	default:
		c.JSON(INVALID_PARAMS.Int(), Response{
			Code:    INVALID_PARAMS,
			Msg:     INVALID_PARAMS.GetMsg(),
			Package: "Unknown chaincode language",
		})
		return
	}

	err = os.MkdirAll(dst, 0777)
	if err != nil {
		c.JSON(401, Response{
			Code:    401,
			Msg:     "Failed to creat chaincode saving dir",
			Package: "Failed to creat chaincode saving dir",
		})
		return
	}

	//存储文件
	for _, file := range files {
		err = c.SaveUploadedFile(file, dst+string(os.PathSeparator)+file.Filename)
		if err != nil {
			c.JSON(ERROR.Int(), Response{
				Code:    ERROR,
				Msg:     ERROR.GetMsg(),
				Package: err.Error(),
			})
			return
		}
	}

	//创建合约

	conf := contract.ChaincodeConfig{
		ChannelName:      chName,
		ChaincodeName:    ccName,
		ChaincodeDesc:    ccDesc,
		ChaincodeGoPath:  utils.ConfigPathWithId(bid) + "/chaincode/go/",
		ChaincodePath:    ccpath,
		ChaincodeVersion: ccVer,
		ChaincodeLang:    cclang,
	}

	id, err := b.CreatContract(conf)
	if err != nil {
		//合约创建失败，删除上传的合约文件
		_ = os.RemoveAll(dst)
		//返回错误
		c.JSON(http.StatusInternalServerError, Response{
			Code:    BLOCKCHAIN_CONTRACT_ERROR,
			Msg:     BLOCKCHAIN_CONTRACT_ERROR.GetMsg(),
			Package: err.Error(),
		})
		return
	}

	installCC(id, b, c)
}

func ManageContract(c *gin.Context) {
	var config core.ContractManageInfoToJson
	err := c.BindJSON(&config)
	if err != nil {
		c.JSON(INVALID_PARAMS.Int(), Response{
			Code:    INVALID_PARAMS,
			Msg:     INVALID_PARAMS.GetMsg(),
			Package: err.Error(),
		})
		return
	}

	b, ok := core.CoreBlockChainManager.GetBlockchainById(config.BlockchainID)
	if !ok {
		c.JSON(http.StatusInternalServerError, Response{
			Code:    RESOURCES_NOT_FOUND,
			Msg:     RESOURCES_NOT_FOUND.GetMsg(),
			Package: "找不到请求的id对应的区块连网络",
		})
		return
	}

	args := strings.Fields(config.Args)

	switch config.Operation {
	case core.Install:
		installCC(config.ID, b, c)
	case core.Instantiate:
		instantiateCC(config.ID, b, args, c)
	case core.Upgrade:
		upgradeCC(config.ID, b, args, c)
	default:
		c.JSON(INVALID_PARAMS.Int(), Response{
			Code:    INVALID_PARAMS,
			Msg:     INVALID_PARAMS.GetMsg(),
			Package: "未知操作类型",
		})
	}

}

func InvokeContract(c *gin.Context) {
	var config core.ContractInvokeInfoToJson
	err := c.BindJSON(&config)
	if err != nil {
		c.JSON(INVALID_PARAMS.Int(), Response{
			Code:    INVALID_PARAMS,
			Msg:     INVALID_PARAMS.GetMsg(),
			Package: err.Error(),
		})
		return
	}

	b, ok := core.CoreBlockChainManager.GetBlockchainById(config.BlockchainID)
	if !ok {
		c.JSON(http.StatusInternalServerError, Response{
			Code:    RESOURCES_NOT_FOUND,
			Msg:     RESOURCES_NOT_FOUND.GetMsg(),
			Package: "找不到请求的id对应的区块连网络",
		})
		return
	}

	var isQuery bool
	switch config.InvokeType {
	case core.Query:
		isQuery = true
	case core.Invoke:
		isQuery = false
	default:
		c.JSON(INVALID_PARAMS.Int(), Response{
			Code:    INVALID_PARAMS,
			Msg:     INVALID_PARAMS.GetMsg(),
			Package: "合约调用类型必须为invoke或query",
		})
		return
	}

	args := strings.Fields(config.Args)

	result, err := b.InvokeChaincode(config.ID, args, isQuery)
	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			Code:    BLOCKCHAIN_CONTRACT_ERROR,
			Msg:     BLOCKCHAIN_CONTRACT_ERROR.GetMsg(),
			Package: err.Error(),
		})
		return
	}

	c.JSON(SUCCESS.Int(), Response{
		Code:    SUCCESS,
		Msg:     SUCCESS.GetMsg(),
		Package: result,
	})
}

func GetContractLogs(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(INVALID_PARAMS.Int(), Response{
			Code:    INVALID_PARAMS,
			Msg:     INVALID_PARAMS.GetMsg(),
			Package: err.Error(),
		})
		return
	}

	b, ok := core.CoreBlockChainManager.GetBlockchainById(id)
	if !ok {
		c.JSON(http.StatusInternalServerError, Response{
			Code:    RESOURCES_NOT_FOUND,
			Msg:     RESOURCES_NOT_FOUND.GetMsg(),
			Package: "找不到请求的id对应的区块连网络",
		})
		return
	}

	result := b.GetContractLogs()

	c.JSON(SUCCESS.Int(), Response{
		Code:    SUCCESS,
		Msg:     SUCCESS.GetMsg(),
		Package: result,
	})
}

func GetAllContractLogs(c *gin.Context) {
	var result []logger.ContractInvokeLog
	for _, b := range core.CoreBlockChainManager.GetBlockchains() {
		result = mergeContractLogsWithTimeOrder(result, b.GetContractLogs())
	}

	c.JSON(SUCCESS.Int(), Response{
		Code:    SUCCESS,
		Msg:     SUCCESS.GetMsg(),
		Package: result,
	})
}

//--------------------
//此处开始为handler辅助函数

func getBlockchainToJsonFromBlockchain(b *core.Blockchain) core.BlockchainToJson {
	result := core.BlockchainToJson{
		ID:         b.GetId(),
		Name:       b.GetName(),
		OrdererOrg: getOrdererOrgToJsonFromBlockchain(b),
		PeerOrg:    getPeerOrgToJsonsFromBlockchain(b),
		Config:     nil,
		Channels:   getChannelToJsonsFromBlockchain(b),
		Status:     b.GetStatus().String(),
		Netname:    b.GetDockerNetName(),
	}

	return result
}

func getOrdererOrgToJsonFromBlockchain(b *core.Blockchain) core.OrdererOrgToJson {
	conf := b.GetConfigurator()

	if conf.Ordererorgconf == nil {
		return core.OrdererOrgToJson{}
	}

	result := core.OrdererOrgToJson{
		Name:   conf.Ordererorgconf.Name(),
		Domain: conf.Ordererorgconf.Domain(),
		MSPID:  conf.Ordererorgconf.MSPID(),
	}

	if conf.Ordererorgconf.GetNodes()[0] == nil {
		return result
	}

	result.Orderer = core.NodeToJson{
		Name: conf.Ordererorgconf.GetNodes()[0].Host(),
		Port: conf.Ordererorgconf.GetNodes()[0].Port().Port(),
	}

	return result
}

func getPeerOrgToJsonsFromBlockchain(b *core.Blockchain) []core.PeerOrgToJson {
	conf := b.GetConfigurator()

	var result []core.PeerOrgToJson

	for _, peerOrg := range conf.Peerorgsconf {
		peerOrgToJson := core.PeerOrgToJson{
			Name:   peerOrg.Name(),
			Domain: peerOrg.Domain(),
			MSPID:  peerOrg.MSPID(),
			Peers:  nil,
		}
		//向peer组织中添加peer节点
		for _, peer := range peerOrg.GetNodes() {
			pToJson := core.NodeToJson{
				Name: peer.Host(),
				Port: peer.Port().Port(),
			}

			peerOrgToJson.Peers = append(peerOrgToJson.Peers, pToJson)
		}

		result = append(result, peerOrgToJson)
	}

	return result
}

func getChannelToJsonsFromBlockchain(b *core.Blockchain) []core.ChannelToJson {
	var result []core.ChannelToJson

	for _, ch := range b.GetChannelManager().GetChannelinfos() {
		chToJ := core.ChannelToJson{
			Name:           ch.GetChannelName(),
			Peers:          ch.GetPeers(),
			AnchorPeers:    nil,
			BlockchainID:   b.GetId(),
			Blockchainname: b.GetName(),
		}

		result = append(result, chToJ)
	}

	return result
}

func getContractToJsonsFromBlockchain(b *core.Blockchain) []core.ContractToJson {
	var result []core.ContractToJson

	for _, cc := range b.GetContractManager().GetChaincodes() {
		ccToJ := core.ContractToJson{
			ID:              cc.ChaincodeID(),
			BlockchainID:    b.GetId(),
			BlockchainName:  b.GetName(),
			ContractName:    cc.ChaincodeName(),
			ChannelName:     cc.ChannelName(),
			ContractDesc:    cc.ChaincodeDesc(),
			ContractVersion: cc.ChaincodeVer(),
			ContractPath:    "",
			ContractLang:    cc.ChaincodeLanguage(),
		}

		result = append(result, ccToJ)
	}

	return result
}

func installCC(id int, b *core.Blockchain, c *gin.Context) {
	//下面是安装链码的代码段
	resp, err := b.InstallContract(id)
	if err != nil {
		//没有节点安装到链码，完全失败
		if resp == nil {
			c.JSON(http.StatusInternalServerError, Response{
				Code:    BLOCKCHAIN_CONTRACT_ERROR,
				Msg:     BLOCKCHAIN_CONTRACT_ERROR.GetMsg(),
				Package: err.Error(),
			})
			return
		} else { //部分节点安装到链码，部分节点安装失败
			c.JSON(http.StatusInternalServerError, Response{
				Code:    PARTIAL_CONTENT,
				Msg:     PARTIAL_CONTENT.GetMsg(),
				Package: err.Error(),
			})
			return
		}
	}
	c.JSON(SUCCESS.Int(), Response{
		Code:    SUCCESS,
		Msg:     SUCCESS.GetMsg(),
		Package: resp,
	})
}

func instantiateCC(id int, b *core.Blockchain, args []string, c *gin.Context) {
	err := b.InstantiateCC(id, args, false)
	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			Code:    BLOCKCHAIN_CONTRACT_ERROR,
			Msg:     BLOCKCHAIN_CONTRACT_ERROR.GetMsg(),
			Package: err.Error(),
		})
		return
	}

	c.JSON(SUCCESS.Int(), Response{
		Code:    SUCCESS,
		Msg:     SUCCESS.GetMsg(),
		Package: nil,
	})
}

func upgradeCC(id int, b *core.Blockchain, args []string, c *gin.Context) {
	err := b.InstantiateCC(id, args, true)
	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			Code:    BLOCKCHAIN_CONTRACT_ERROR,
			Msg:     BLOCKCHAIN_CONTRACT_ERROR.GetMsg(),
			Package: err.Error(),
		})
		return
	}

	c.JSON(SUCCESS.Int(), Response{
		Code:    SUCCESS,
		Msg:     SUCCESS.GetMsg(),
		Package: nil,
	})
}

func addPeers(name string, b *core.Blockchain, peers []string, c *gin.Context) {
	err := b.AddPeerToChannel(name, peers)
	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			Code:    BLOCKCHAIN_CHANNEL_ERROR,
			Msg:     BLOCKCHAIN_CHANNEL_ERROR.GetMsg(),
			Package: err.Error(),
		})
		return
	}

	c.JSON(SUCCESS.Int(), Response{
		Code: SUCCESS,
		Msg:  SUCCESS.GetMsg(),
	})
}

func mergeContractLogsWithTimeOrder(a []logger.ContractInvokeLog, b []logger.ContractInvokeLog) []logger.ContractInvokeLog {
	result := make([]logger.ContractInvokeLog, len(a)+len(b))
	var i, j, k = 0, 0, 0
	for i < len(a) && j < len(b) {
		if a[i].Time.Before(b[j].Time) {
			result[k] = a[i]
			i++
			k++
		} else {
			result[k] = b[j]
			j++
			k++
		}
	}

	for i < len(a) {
		result[k] = a[i]
		i++
		k++
	}

	for j < len(b) {
		result[k] = b[j]
		j++
		k++
	}

	return result

}
