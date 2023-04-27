package deploy

import (
	"errors"
	"os"
)

type PeerOrgConfig struct {
	nodes      []NodeConfigIf
	name       string
	domain     string
	mspid      string
	anchor     NodeConfigIf
	configpath string
}

func (p *PeerOrgConfig) Name() string {
	return p.name
}

func (p *PeerOrgConfig) Domain() string {
	return p.domain
}

func (p *PeerOrgConfig) MSPID() string {
	return p.mspid
}

func (p *PeerOrgConfig) CryptoPath() string {
	return p.configpath + string(os.PathSeparator) + "crypto-config" + string(os.PathSeparator) + "peerOrganizations" + string(os.PathSeparator) + p.domain
}

func (p *PeerOrgConfig) ConfigtxPath() string {
	return p.configpath + string(os.PathSeparator) + "configtx"
}

func (p *PeerOrgConfig) AddNode(node NodeConfigIf) error {
	for _, n := range p.nodes {
		if n.Host() == node.Host() {
			return errors.New("node config already exists !")
		}
	}
	p.nodes = append(p.nodes, node)
	return nil
}

func (p *PeerOrgConfig) DelNode(host string) error {
	for i, n := range p.nodes {
		if n.Host() == host {
			p.nodes = append(p.nodes[:i], p.nodes[i+1:]...)
			return nil
		}
	}
	return errors.New("node config does not exists !")
}

func (p *PeerOrgConfig) GetAnchorPeer() NodeConfigIf {
	return p.anchor
}

func (p *PeerOrgConfig) GetNodes() []NodeConfigIf {
	return p.nodes
}

func NewPeerOrgConfig(name string, domain string, mspid string, path string) *PeerOrgConfig {
	return &PeerOrgConfig{name: name, domain: domain, mspid: mspid, configpath: path}
}
