package core

type BlockchainStatus int

const (
	Configuring BlockchainStatus = iota
	Running
	Stop
)

var BlockchainStatusMap = map[BlockchainStatus]string{
	Configuring: "Configuring",
	Running:     "Running",
	Stop:        "Stop",
}

func (s BlockchainStatus) String() string {
	if result, ok := BlockchainStatusMap[s]; ok {
		return result
	} else {
		return "Unknown status"
	}
}

func (s BlockchainStatus) Valid() bool {
	_, ok := BlockchainStatusMap[s]
	return ok
}
