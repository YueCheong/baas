package network

import (
	"bit-bass/deploy"
	"context"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"io"
	"os"
	"strconv"
)

type GossipConfigIf interface {
	// gossip external endpoint
	EndPoint() string
	// gossip bootstrap
	SetBootAddr(addr string)
}

type PeerNode struct {
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
	// chaincode port
	ccport nat.Port
	// gossip bootstrap
	bootaddr string
	// crypto materials path
	cryptopath string
	// configtx materails path
	configtxpath string
	// 节点连接的网络名
	dockernetname string
}

func (p *PeerNode) EndPoint() string {
	return p.ContainerName() + ":" + p.ServePort().Port()
}

func (p *PeerNode) SetBootAddr(addr string) {
	p.bootaddr = addr
}

func (p *PeerNode) Create() error {
	var binds []string
	binds = append(binds, "/var/run/:/host/var/run/")
	binds = append(binds, p.configtxpath+":/etc/hyperledger/configtx")
	binds = append(binds, p.cryptopath+string(os.PathSeparator)+"peers"+string(os.PathSeparator)+p.ContainerName()+string(os.PathSeparator)+"msp:/etc/hyperledger/msp/peer")
	binds = append(binds, p.cryptopath+string(os.PathSeparator)+"users:/etc/hyperledger/msp/users")

	resp, err := p.cli.ContainerCreate(p.ctx,
		&container.Config{
			Image: PEERIMAGE,
			Env: []string{
				"FABRIC_LOGGING_SPEC=INFO",
				"CORE_CHAINCODE_LOGGING_LEVEL=info",
				"CORE_VM_ENDPOINT=unix:///host/var/run/docker.sock",
				"CORE_PEER_ID=" + p.ContainerName(),
				"CORE_PEER_ADDRESS=" + p.ContainerName() + ":" + p.port.Port(),
				"CORE_PEER_LISTENADDRESS=0.0.0.0:" + p.port.Port(),
				"CORE_PEER_CHAINCODEADDRESS=" + p.ContainerName() + ":" + p.ccport.Port(),
				"CORE_PEER_CHAINCODELISTENADDRESS=0.0.0.0:" + p.ccport.Port(),
				"CORE_PEER_GOSSIP_BOOTSTRAP=" + p.bootaddr,
				"CORE_PEER_GOSSIP_EXTERNALENDPOINT=" + p.EndPoint(),
				"CORE_PEER_LOCALMSPID=" + p.org.MSPID(),
				"CORE_PEER_MSPCONFIGPATH=/etc/hyperledger/msp/peer/",
				"CORE_VM_DOCKER_HOSTCONFIG_NETWORKMODE=" + p.dockernetname,
			},
			Cmd: []string{"peer", "node", "start"},
			//WorkingDir: "/opt/gopath/src/github.com/hyperledger/fabric/peer",
			WorkingDir: "/etc/hyperledger/configtx",
			ExposedPorts: nat.PortSet{
				p.ServePort(): struct{}{},
				p.ccport:      struct{}{},
			},
		},
		&container.HostConfig{
			Binds: binds,
			PortBindings: nat.PortMap{
				p.ServePort(): []nat.PortBinding{{HostIP: "0.0.0.0", HostPort: p.ServePort().Port()}},
				p.ccport:      []nat.PortBinding{{HostIP: "0.0.0.0", HostPort: p.ccport.Port()}},
			},
		},
		&network.NetworkingConfig{},
		p.ContainerName())
	if err != nil {
		return err
	}
	p.cid = resp.ID
	return nil
}

func (p *PeerNode) Remove() error {
	return p.cli.ContainerRemove(p.ctx, p.cid, types.ContainerRemoveOptions{})
}

func (p *PeerNode) Start() error {
	return p.cli.ContainerStart(p.ctx, p.cid, types.ContainerStartOptions{})
}

func (p *PeerNode) Stop() error {
	return p.cli.ContainerStop(p.ctx, p.cid, nil)
}

func (p *PeerNode) ConnectNet(net deploy.NetIf) error {
	return p.cli.NetworkConnect(p.ctx, net.NetID(), p.cid, &network.EndpointSettings{})
}

func (p *PeerNode) NodeName() string {
	return p.name
}

func (p *PeerNode) NodeOrg() OrgIf {
	return p.org
}

func (p *PeerNode) Inspect() (types.ContainerJSON, error) {
	return p.cli.ContainerInspect(p.ctx, p.cid)
}

func (p *PeerNode) PrintLog() (int64, error) {
	out, err := p.cli.ContainerLogs(p.ctx, p.cid, types.ContainerLogsOptions{ShowStdout: true, ShowStderr: true})
	if err != nil {
		return 0, err
	}
	return io.Copy(os.Stdout, out)
}

func (p *PeerNode) ContainerID() string {
	return p.cid
}

func (p *PeerNode) ContainerName() string {
	return p.name + "." + p.org.Domain()
}

func (p *PeerNode) ServePort() nat.Port {
	return p.port
}

func NewPeerNode(node deploy.NodeConfigIf) *PeerNode {
	ctx := context.Background()
	cli, err := client.NewEnvClient()
	if err != nil {
		panic(err)
	}

	name := node.Host()
	port := node.Port()
	ccport, err := nat.NewPort(node.Port().Proto(), strconv.Itoa(node.Port().Int()+1))
	cryptopath := node.Domain().CryptoPath()
	configtxpath := node.Domain().ConfigtxPath()

	pn := &PeerNode{
		ctx:           ctx,
		cli:           cli,
		name:          name,
		port:          port,
		ccport:        ccport,
		cryptopath:    cryptopath,
		configtxpath:  configtxpath,
		dockernetname: node.NetName(),
	}

	return pn
}
