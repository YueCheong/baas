package network

import (
	"bit-bass/deploy"
	"errors"
)

type PeerOrg struct {
	name string
	// 10086.cn, bit.edu.cn
	domain string
	// Org86MSP, OrgBitMSP
	mspid string
	//
	peers []*PeerNode
}

func (po *PeerOrg) Create() error {
	for _, peer := range po.peers {
		if err := peer.Create(); err != nil {
			return err
		}
	}
	return nil
}

func (po *PeerOrg) Remove() error {
	for _, peer := range po.peers {
		if err := peer.Remove(); err != nil {
			return err
		}
	}
	return nil
}

func (po *PeerOrg) Start() error {
	for _, peer := range po.peers {
		if err := peer.Start(); err != nil {
			return err
		}
	}
	return nil
}

func (po *PeerOrg) Stop() error {
	for _, peer := range po.peers {
		if err := peer.Stop(); err != nil {
			return err
		}
	}
	return nil
}

func (po *PeerOrg) ConnectNet(net deploy.NetIf) error {
	for _, peer := range po.peers {
		if err := peer.ConnectNet(net); err != nil {
			return err
		}
	}
	return nil
}

func (po *PeerOrg) Name() string {
	return po.name
}

func (po *PeerOrg) Domain() string {
	return po.domain
}

func (po *PeerOrg) MSPID() string {
	return po.mspid
}

func (po *PeerOrg) AddPeerNode(node *PeerNode) error {
	for _, p := range po.peers {
		if p.NodeName() == node.NodeName() {
			return errors.New("The peer is already exists !")
		}
	}
	node.org = po
	po.peers = append(po.peers, node)
	return nil
}

func (po *PeerOrg) DelPeerNode(nodeid string) error {
	for i, p := range po.peers {
		if p.NodeName() == nodeid {
			po.peers = append(po.peers[:i], po.peers[i+1:]...)
			return nil
		}
	}
	return errors.New("node config does not exists !")
}

func (po *PeerOrg) GetPeerNodes() []*PeerNode {
	return po.peers
}

func NewPeerOrg(org deploy.OrgConfigIf) *PeerOrg {

	peerorg := &PeerOrg{
		name:   org.Name(),
		domain: org.Domain(),
		mspid:  org.MSPID(),
	}

	for _, node := range org.GetNodes() {
		_ = peerorg.AddPeerNode(NewPeerNode(node))
	}

	peers := peerorg.GetPeerNodes()
	count := len(peers)
	if count > 0 {
		for i := 0; i < count-1; i++ {
			peers[i].SetBootAddr(peers[i+1].EndPoint())
		}
		peers[count-1].SetBootAddr(peers[0].EndPoint())
	}

	return peerorg
}
