package channel

import (
	"bit-bass/network"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
	"github.com/pkg/errors"
)

type Channelinfo struct {
	channelName string
	peers       []string
	created     bool
	configPath  string
}

func (c *Channelinfo) GetChannelName() string {
	return c.channelName
}

func (c *Channelinfo) GetPeers() []string {
	return c.peers
}

func (c *Channelinfo) NewChannelClient(user string, org network.OrgIf) (*ChannelClient, error) {
	netConifg := config.FromFile(c.configPath + "/sdk-config.yaml")
	fabricsdk, err := fabsdk.New(netConifg)
	if err != nil {
		return nil, errors.Errorf("Faild to creat fabric sdk : %v", err)
	}

	channel := ChannelClient{
		info:        c,
		operateOrg:  org,
		operateUser: user,
		sdk:         fabricsdk,
	}
	return &channel, nil
}
