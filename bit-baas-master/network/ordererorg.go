package network

import (
	"bit-bass/deploy"
	"errors"
)

type OrdererOrg struct {
	// orderer
	name string
	// example.com
	domain string
	// OrdererMSP
	mspid string
	// 排序节点
	orderer *OrdererNode
}

func (oo *OrdererOrg) Create() error {
	return oo.orderer.Create()
}

func (oo *OrdererOrg) Remove() error {
	return oo.orderer.Remove()
}

func (oo *OrdererOrg) Start() error {
	return oo.orderer.Start()
}

func (oo *OrdererOrg) Stop() error {
	return oo.orderer.Stop()
}

func (oo *OrdererOrg) ConnectNet(net deploy.NetIf) error {
	return oo.orderer.ConnectNet(net)
}

func (oo *OrdererOrg) Name() string {
	return oo.name
}

func (oo *OrdererOrg) Domain() string {
	return oo.domain
}

func (oo *OrdererOrg) MSPID() string {
	return oo.mspid
}

func (oo *OrdererOrg) AddOrdererNode(node *OrdererNode) error {
	if oo.orderer != nil {
		return errors.New("orderer node exists!")
	}
	node.org = oo
	oo.orderer = node
	return nil
}

func (oo *OrdererOrg) DelOrdererNode(name string) error {
	if oo.orderer.NodeName() != name {
		return errors.New("orderer node does not exists!")
	}
	oo.orderer = nil
	return nil
}

func (oo *OrdererOrg) GetOrdererNode() *OrdererNode {
	return oo.orderer
}

func NewOrdererOrg(org deploy.OrgConfigIf) *OrdererOrg {

	ordererorg := &OrdererOrg{
		name:   org.Name(),
		domain: org.Domain(),
		mspid:  org.MSPID(),
	}

	node := org.GetNodes()[0]
	_ = ordererorg.AddOrdererNode(NewOrdererNode(node))

	return ordererorg
}
