package deploy

import "github.com/docker/go-connections/nat"

type NodeConfig struct {
	host          string
	domain        OrgConfigIf
	port          nat.Port
	dockerNetName string
}

func (n *NodeConfig) Host() string {
	return n.host
}

func (n *NodeConfig) Domain() OrgConfigIf {
	return n.domain
}

func (n *NodeConfig) Port() nat.Port {
	return n.port
}

func (n *NodeConfig) JoinOrg(domain OrgConfigIf) error {
	if err := domain.AddNode(n); err != nil {
		return err
	} else {
		n.domain = domain
		return nil
	}
}

func (n *NodeConfig) NetName() string {
	return n.dockerNetName
}

func NewNodeConfig(host string, domain OrgConfigIf, port nat.Port, netname string) *NodeConfig {
	if domain == nil {
		return &NodeConfig{host: host, port: port, dockerNetName: netname}
	} else {
		node := &NodeConfig{host: host, domain: domain, port: port, dockerNetName: netname}
		if err := domain.AddNode(node); err != nil {
			panic(err)
		}
		return node
	}
}
