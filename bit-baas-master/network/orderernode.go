package network

import (
	"bit-bass/deploy"
	"context"
	"errors"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"io"
	"os"
)

type PeerOrgCryptoIf interface {
	// 添加peer组织的信息
	AddOrgCrypto(org deploy.OrgConfigIf) error
	// 删除指定peer组织信息
	DelOrgCrypto(name string) error
	// 获取peer组织的密钥信息
	GetOrgCrypto() []OrgCryptoInfo
}

type OrgCryptoInfo struct {
	orgid   string
	crypath string
}

type OrdererNode struct {
	ctx context.Context
	cli *client.Client
	// container_id
	cid string
	// peer0, peer1
	name string
	// org interface
	org OrgIf
	// serve port
	port nat.Port
	// crypto materials path
	cryptopath string
	// configtx materails path
	configtxpath string
	// peers orgnization crypto materials path
	peersinfo []OrgCryptoInfo
	// 节点链接的docker网络名
	dockernetname string
}

func (o *OrdererNode) AddOrgCrypto(org deploy.OrgConfigIf) error {
	for _, pi := range o.peersinfo {
		if pi.orgid == "peer"+org.Name() {
			return errors.New("org crypto exsit!")
		}
	}
	info := OrgCryptoInfo{"peer" + org.Name(), org.CryptoPath() + string(os.PathSeparator) + "peers" + string(os.PathSeparator) + org.GetNodes()[0].Host() + "." + org.Domain()}
	o.peersinfo = append(o.peersinfo, info)
	return nil
}

func (o *OrdererNode) DelOrgCrypto(name string) error {
	for i, pi := range o.peersinfo {
		if pi.orgid == "peer"+name {
			o.peersinfo = append(o.peersinfo[:i], o.peersinfo[i+1:]...)
			return nil
		}
	}
	return errors.New("peer crypto does not exists !")
}

func (o *OrdererNode) GetOrgCrypto() []OrgCryptoInfo {
	return o.peersinfo
}

func (o *OrdererNode) Create() error {
	var binds []string
	binds = append(binds, o.configtxpath+":/etc/hyperledger/configtx")
	binds = append(binds, o.cryptopath+string(os.PathSeparator)+"orderers"+string(os.PathSeparator)+o.ContainerName()+":/etc/hyperledger/msp/orderer")
	for _, pi := range o.peersinfo {
		binds = append(binds, pi.crypath+":/etc/hyperledger/msp"+string(os.PathSeparator)+pi.orgid)
	}

	resp, err := o.cli.ContainerCreate(o.ctx,
		&container.Config{
			Image: ORDERERIMAGE,
			Env: []string{
				"FABRIC_LOGGING_SPEC=INFO",
				"ORDERER_GENERAL_LISTENADDRESS=0.0.0.0",
				"ORDERER_GENERAL_GENESISMETHOD=file",
				"ORDERER_GENERAL_GENESISFILE=/etc/hyperledger/configtx/genesis.block",
				"ORDERER_GENERAL_LOCALMSPID=" + o.org.MSPID(),
				"ORDERER_GENERAL_LOCALMSPDIR=/etc/hyperledger/msp/orderer/msp",
				"CORE_VM_DOCKER_HOSTCONFIG_NETWORKMODE=" + o.dockernetname,
				"CORE_VM_ENDPOINT=unix:///host/var/run/docker.sock",
				"ORDERER_GENERAL_LISTENPORT=" + o.ServePort().Port(),
			},
			Cmd:        []string{"orderer"},
			WorkingDir: "/opt/gopath/src/github.com/hyperledger/fabric",
			ExposedPorts: nat.PortSet{
				o.ServePort(): struct{}{},
			},
		},
		&container.HostConfig{
			Binds: binds,
			PortBindings: nat.PortMap{
				o.ServePort(): []nat.PortBinding{{HostIP: "0.0.0.0", HostPort: o.ServePort().Port()}},
			},
		},
		&network.NetworkingConfig{},
		o.ContainerName())
	if err != nil {
		return err
	}
	o.cid = resp.ID
	return nil
}

func (o OrdererNode) Remove() error {
	return o.cli.ContainerRemove(o.ctx, o.cid, types.ContainerRemoveOptions{})
}

func (o OrdererNode) Start() error {
	return o.cli.ContainerStart(o.ctx, o.cid, types.ContainerStartOptions{})
}

func (o OrdererNode) Stop() error {
	return o.cli.ContainerStop(o.ctx, o.cid, nil)
}

func (o OrdererNode) ConnectNet(net deploy.NetIf) error {
	return o.cli.NetworkConnect(o.ctx, net.NetID(), o.cid, &network.EndpointSettings{})
}

func (o OrdererNode) NodeName() string {
	return o.name
}

func (o OrdererNode) NodeOrg() OrgIf {
	return o.org
}

func (o OrdererNode) Inspect() (types.ContainerJSON, error) {
	return o.cli.ContainerInspect(o.ctx, o.cid)
}

func (o OrdererNode) PrintLog() (int64, error) {
	out, err := o.cli.ContainerLogs(o.ctx, o.cid, types.ContainerLogsOptions{ShowStdout: true, ShowStderr: true})
	if err != nil {
		return 0, err
	}
	return io.Copy(os.Stdout, out)
}

func (o OrdererNode) ContainerID() string {
	return o.cid
}

func (o OrdererNode) ContainerName() string {
	return o.name + "." + o.org.Domain()
}

func (o OrdererNode) ServePort() nat.Port {
	return o.port
}

func NewOrdererNode(node deploy.NodeConfigIf) *OrdererNode {
	ctx := context.Background()
	cli, err := client.NewEnvClient()
	if err != nil {
		panic(err)
	}

	name := node.Host()
	port := node.Port()

	cryptopath := node.Domain().CryptoPath()
	configtxpath := node.Domain().ConfigtxPath()

	on := &OrdererNode{
		ctx:           ctx,
		cli:           cli,
		name:          name,
		port:          port,
		cryptopath:    cryptopath,
		configtxpath:  configtxpath,
		dockernetname: node.NetName(),
	}

	return on
}
