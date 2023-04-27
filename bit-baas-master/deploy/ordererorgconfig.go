package deploy

import (
	"errors"
	"os"
)

type OrdererOrgConfig struct {
	orderer    NodeConfigIf
	name       string
	domain     string
	mspid      string
	configpath string
}

func (o *OrdererOrgConfig) Name() string {
	return o.name
}

func (o *OrdererOrgConfig) Domain() string {
	return o.domain
}

func (o *OrdererOrgConfig) MSPID() string {
	return o.mspid
}

func (o *OrdererOrgConfig) CryptoPath() string {
	return o.configpath + string(os.PathSeparator) + "crypto-config" + string(os.PathSeparator) + "ordererOrganizations" + string(os.PathSeparator) + o.domain
}

func (o *OrdererOrgConfig) ConfigtxPath() string {
	return o.configpath + string(os.PathSeparator) + "configtx"
}

func (o *OrdererOrgConfig) AddNode(node NodeConfigIf) error {
	o.orderer = node
	return nil
}

func (o *OrdererOrgConfig) DelNode(host string) error {
	if o.orderer.Host() != host {
		return errors.New("orderer config does not exists !")
	}
	o.orderer = nil
	return nil
}

func (o *OrdererOrgConfig) GetNodes() []NodeConfigIf {
	return []NodeConfigIf{o.orderer}

}

func NewOrdererOrgConfig(name string, domain string, mspid string, path string) *OrdererOrgConfig {
	return &OrdererOrgConfig{name: name, domain: domain, mspid: mspid, configpath: path}
}
