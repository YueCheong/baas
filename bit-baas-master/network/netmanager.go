package network

import (
	"bit-bass/deploy"
	"errors"
)

type NetManager struct {
	Netop      deploy.NetIf // 网络接口
	Ordererorg *OrdererOrg  // orderer组织
	Peerorgs   []*PeerOrg   // peer组织
	Client     *Tools
}

func (n *NetManager) LoadConfig(config *deploy.Configurator) error {
	//用网络路径中的最后一部分（也就是ID）来标识不同的网络
	if config.DockerNet == nil {
		return errors.New("Network in configurator is nil")
	}
	n.Netop = config.DockerNet
	n.Client = NewTools(config.Cli_name, config.Cli_node)
	n.Ordererorg = NewOrdererOrg(config.Ordererorgconf)
	for _, posc := range config.Peerorgsconf {
		n.Peerorgs = append(n.Peerorgs, NewPeerOrg(posc))
		if err := n.Ordererorg.orderer.AddOrgCrypto(posc); err != nil {
			return err
		}
	}
	return nil
}

func (n *NetManager) Prepared() error {
	if !n.Netop.IsNetExist() {
		if err := n.Netop.CreateNet(); err != nil {
			return err
		}
	}
	if err := n.Client.Create(); err != nil {
		return nil
	}
	if err := n.Client.ConnectNet(n.Netop); err != nil {
		return err
	}
	if err := n.Ordererorg.Create(); err != nil {
		return err
	}
	if err := n.Ordererorg.ConnectNet(n.Netop); err != nil {
		return err
	}
	for _, po := range n.Peerorgs {
		if err := po.Create(); err != nil {
			return err
		}
		if err := po.ConnectNet(n.Netop); err != nil {
			return err
		}
	}
	return nil
}

func (n *NetManager) StartNet() error {
	if err := n.Client.Start(); err != nil {
		return nil
	}
	if err := n.Ordererorg.Start(); err != nil {
		return err
	}
	for _, po := range n.Peerorgs {
		if err := po.Start(); err != nil {
			return err
		}
	}
	return nil
}

func (n *NetManager) StopNet() error {
	if err := n.Client.Stop(); err != nil {
		return nil
	}
	if err := n.Ordererorg.Stop(); err != nil {
		return err
	}
	for _, po := range n.Peerorgs {
		if err := po.Stop(); err != nil {
			return err
		}
	}
	return nil
}

func (n *NetManager) Remove() error {
	if err := n.Client.Remove(); err != nil {
		return nil
	}
	if err := n.Ordererorg.Remove(); err != nil {
		return err
	}
	for _, po := range n.Peerorgs {
		if err := po.Remove(); err != nil {
			return err
		}
	}
	return nil
}

func NewNetManager() *NetManager {
	return &NetManager{}
}
