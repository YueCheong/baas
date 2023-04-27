package deploy

import (
	"errors"
	"github.com/docker/go-connections/nat"
)

type Configurator struct {
	DockerNet      NetIf
	Outpath        string
	Ordererorgconf *OrdererOrgConfig
	Peerorgsconf   []*PeerOrgConfig
	Cli_name       string
	Cli_node       NodeConfigIf
}

func (c *Configurator) SetNetwork(netIf NetIf) error {

	return errors.New("Net didn't exist")
}

func (c *Configurator) SetOrdererOrg(name string, domain string, mspid string) error {
	c.Ordererorgconf = NewOrdererOrgConfig(name, domain, mspid, c.Outpath)
	return nil
}

func (c *Configurator) SetOrdererNode(host string, port nat.Port) error {
	if c.Ordererorgconf != nil {
		NewNodeConfig(host, c.Ordererorgconf, port, c.DockerNet.NetName())
		return nil
	}
	return errors.New("setting a orderer org first!")
}

func (c *Configurator) AddPeerOrg(name string, domain string, mspid string) error {
	for _, pc := range c.Peerorgsconf {
		if pc.Name() == name {
			return errors.New("peer org config already exists !")
		}
	}
	psc := NewPeerOrgConfig(name, domain, mspid, c.Outpath)
	c.Peerorgsconf = append(c.Peerorgsconf, psc)
	return nil
}

func (c *Configurator) DelPeerOrg(name string) error {
	for i, pc := range c.Peerorgsconf {
		if pc.Name() == name {
			c.Peerorgsconf = append(c.Peerorgsconf[:i], c.Peerorgsconf[i+1:]...)
			return nil
		}
	}
	return errors.New("peer config does not exists !")
}

func (c *Configurator) AddPeerNodeToOrg(orgname string, host string, port nat.Port) error {
	for _, pc := range c.Peerorgsconf {
		if pc.Name() == orgname {
			NewNodeConfig(host, pc, port, c.DockerNet.NetName())
			return nil
		}
	}
	return errors.New("Peer org must exists before add node!")
}

func (c *Configurator) SetAnchorPeerToOrg(orgname string, host string) error {
	for _, pc := range c.Peerorgsconf {
		if pc.Name() == orgname {
			for _, peer := range pc.GetNodes() {
				if peer.Host() == host {
					pc.anchor = peer
					return nil
				}
			}
		}
	}
	return errors.New("Org or peer didn't exist!")
}

func (c *Configurator) DelPeerNodeFromOrg(orgname string, nodehost string) error {
	for _, pc := range c.Peerorgsconf {
		if pc.Name() == orgname {
			return pc.DelNode(nodehost)
		}
	}
	return errors.New("Peer org must exists before del node!")
}

func (c *Configurator) SetToolsCli(name string, node NodeConfigIf) error {
	c.Cli_name = name
	c.Cli_node = node
	return nil
}

func NewConfigurator(configPath string, net NetIf) *Configurator {
	if net.IsNetExist() {
		return &Configurator{
			Outpath:   configPath,
			DockerNet: net,
		}
	}
	return nil
}
