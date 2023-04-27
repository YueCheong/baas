package logger

import (
	"bit-bass/contract"
	"bit-bass/utils"
	"encoding/json"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"time"
)

type Time time.Time

func (t Time) MarshalJSON() ([]byte, error) {
	return json.Marshal(time.Time(t).Format("2006/01/02 15:04:05"))
}

func (t Time) Before(b Time) bool {
	return time.Time(t).Before(time.Time(b))
}

//一条合约调用记录
type ContractInvokeLog struct {
	//合约调用记录的id
	RecordID int
	//调用合约的信息
	ContractID     int
	BlockchainID   int
	BlockchainName string
	ContractName   string
	ChannelName    string
	ContractDesc   string
	ContractVer    string
	//调用的类型，可以为"Query"或"Invoke"
	InvokeType string
	//合约调用的输入参数
	Args []string
	//合约调用的返回
	StatusCode    int32
	TransactionID string
	Payload       string
	//合约调用的时间
	Time Time
}

type ContractInvokeLogger struct {
	blockchainId   int
	blockchainName string
	idGen          *utils.AutoIncIDGen
	logs           []ContractInvokeLog
}

func NewContractInvokeLogger(id int, name string) *ContractInvokeLogger {
	logger := ContractInvokeLogger{}

	logger.idGen = utils.NewAutoIncID()
	logger.blockchainId = id
	logger.blockchainName = name
	return &logger
}

func (l *ContractInvokeLogger) Log(info contract.ChaincodeInfo, args []string, resp channel.Response, query bool) {
	log := ContractInvokeLog{
		RecordID:       l.idGen.GenID(),
		ContractID:     info.ChaincodeID(),
		BlockchainID:   l.blockchainId,
		BlockchainName: l.blockchainName,
		ContractName:   info.ChaincodeName(),
		ChannelName:    info.ChannelName(),
		ContractDesc:   info.ChaincodeDesc(),
		ContractVer:    info.ChaincodeVer(),
		Args:           args,
		StatusCode:     resp.ChaincodeStatus,
		TransactionID:  string(resp.TransactionID),
		Payload:        string(resp.Payload),
		Time:           Time(time.Now()),
	}

	if query {
		log.InvokeType = "Query"
	} else {
		log.InvokeType = "Invoke"
	}

	l.logs = append(l.logs, log)
}

func (l *ContractInvokeLogger) GetLog() []ContractInvokeLog {
	return l.logs
}
