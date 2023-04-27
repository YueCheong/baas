package network

import (
	"bit-bass/deploy"
	"bit-bass/utils"
	"context"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"os"
)

type Tools struct {
	ctx context.Context
	cli *client.Client
	// container_id
	cid string
	// cli
	name string
	// connected peer node
	peer_host       string
	peer_port       string
	peer_domain     string
	peer_localmspid string
	// crypto materials path
	cryptopath string
	// configtx materails path
	configtxpath string
	// chaincode path
	chaincodepath string
}

func (t *Tools) Create() error {
	var binds []string //Add volume
	binds = append(binds, "/var/run/:/host/var/run/")
	binds = append(binds, t.chaincodepath+":/opt/gopath/src/github.com/chaincode")
	binds = append(binds, t.cryptopath+":/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/")
	binds = append(binds, t.configtxpath+":/opt/gopath/src/github.com/hyperledger/fabric/peer/channel-artifacts")

	resp, err := t.cli.ContainerCreate(t.ctx,
		&container.Config{
			Image: TOOLSIMAGE,
			Tty:   true,
			Env: []string{
				"CORE_VM_ENDPOINT=unix:///host/var/run/docker.sock",
				"FABRIC_LOGGING_SPEC=INFO",
				"CORE_PEER_ID=" + t.name,                                                     // cli
				"CORE_PEER_ADDRESS=" + t.peer_host + "." + t.peer_domain + ":" + t.peer_port, //peer0.org1.example.com:7051
				"CORE_PEER_LOCALMSPID=" + t.peer_localmspid,                                  // Org1MSP
				"CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/users/Admin@" + t.peer_domain + "/msp",
			},
			Cmd:        []string{"/bin/bash"},
			WorkingDir: "/opt/gopath/src/github.com/hyperledger/fabric/peer",
		},
		&container.HostConfig{
			Binds: binds,
		},
		&network.NetworkingConfig{},
		"cli"+"."+t.peer_host+"."+t.peer_domain)
	if err != nil {
		return err
	}
	t.cid = resp.ID
	return nil
}

func (t *Tools) Remove() error {
	return t.cli.ContainerRemove(t.ctx, t.cid, types.ContainerRemoveOptions{})
}

func (t *Tools) Start() error {
	return t.cli.ContainerStart(t.ctx, t.cid, types.ContainerStartOptions{})
}

func (t *Tools) Stop() error {
	return t.cli.ContainerStop(t.ctx, t.cid, nil)
}

func (t *Tools) ConnectNet(net deploy.NetIf) error {
	return t.cli.NetworkConnect(t.ctx, net.NetID(), t.cid, &network.EndpointSettings{})
}

//func (tools *Tools) Start() error {
//	ctx := tools.connter.ctx
//	cli := tools.connter.cli
//
//	cPath, _ := os.Getwd()
//	var binds []string //Add volume
//	binds = append(binds, "/var/run/:/host/var/run/")
//	binds = append(binds, cPath+"/../contract/:/opt/gopath/src/github.com/chaincode")
//	binds = append(binds, cPath+"/artifacts/crypto-artifacts:/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/")
//	binds = append(binds, cPath+"/artifacts/configtx:/opt/gopath/src/github.com/hyperledger/fabric/peer/channel-artifacts")
//
//	resp, err := cli.ContainerCreate(ctx,
//		&container.Config{
//			Image: TOOLSIMAGE,
//			Env: []string{
//				"FABRIC_LOGGING_SPEC=INFO",
//				"CORE_PEER_ID=" + tools.ContainerName(),
//				"CORE_PEER_ADDRESS=" + tools.ContainerName() + "+" + tools.port.Port(),
//				"CORE_PEER_LOCALMSPID=" + tools.org.MSPID(),
//				"CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/" +
//					tools.org.Domain() + "/users/" + tools.user + "@" + tools.org.Domain() + "/msp",
//			},
//			Cmd:        []string{"/bin/bash"},
//			WorkingDir: "/opt/gopath/src/github.com/hyperledger/fabric/peer",
//		},
//		&container.HostConfig{
//			Binds: binds,
//			PortBindings: nat.PortMap{
//				tools.ServePort(): []nat.PortBinding{{HostIP: "0.0.0.0", HostPort: tools.ServePort().Port()}},
//			},
//		},
//		&network.NetworkingConfig{},
//		tools.ContainerName())
//
//	if err != nil {
//		return err
//	}
//
//	tools.cid = resp.ID
//	if err := cli.NetworkConnect(ctx, tools.connter.nid, tools.cid, &network.EndpointSettings{}); err != nil {
//		return err
//	}
//	if err := cli.ContainerStart(ctx, tools.cid, types.ContainerStartOptions{}); err != nil {
//		return err
//	}
//
//	return nil
//}
//
//func (tools *Tools) Stop() error {
//	ctx := tools.connter.ctx
//	cli := tools.connter.cli
//
//	return cli.ContainerStop(ctx, tools.cid, nil)
//}
//
//func (tools *Tools) Remove() error {
//	ctx := tools.connter.ctx
//	cli := tools.connter.cli
//
//	return cli.ContainerRemove(ctx, tools.cid, types.ContainerRemoveOptions{})
//}
//
//func (tools *Tools) Inspect() (types.ContainerJSON, error) {
//	ctx := tools.connter.ctx
//	cli := tools.connter.cli
//
//	return cli.ContainerInspect(ctx, tools.cid)
//}
//
//func (tools *Tools) PrintLog() (int64, error) {
//	ctx := tools.connter.ctx
//	cli := tools.connter.cli
//
//	out, err := cli.ContainerLogs(ctx, tools.cid, types.ContainerLogsOptions{ShowStdout: true, ShowStderr: true})
//	if err != nil {
//		return 0, err
//	}
//	return io.Copy(os.Stdout, out)
//}
//
//func (tools *Tools) ContainerID() string {
//	return tools.cid
//}
//
//func (tools *Tools) ContainerName() string {
//	if tools.org.Domain() != "" {
//		return tools.name + "." + tools.org.Domain()
//	} else {
//		return tools.name
//	}
//}
//
//func (tools *Tools) ServePort() nat.Port {
//	return tools.port
//}
//
func NewTools(name string, node deploy.NodeConfigIf) *Tools {
	ctx := context.Background()
	cli, err := client.NewEnvClient()
	if err != nil {
		panic(err)
	}

	peer_host := node.Host()
	peer_port := node.Port().Port()
	peer_domain := node.Domain().Domain()
	peer_localmspid := node.Domain().MSPID()

	cryptopath := node.Domain().CryptoPath()     // utils.ConfigPath() + string(os.PathSeparator) + "crypto-config"
	configtxpath := node.Domain().ConfigtxPath() // utils.ConfigPath() + string(os.PathSeparator) + "configtx"
	chaincodepath := utils.ConfigPath() + string(os.PathSeparator) + "chaincode"

	tools := &Tools{
		ctx:             ctx,
		cli:             cli,
		name:            name,
		peer_host:       peer_host,
		peer_port:       peer_port,
		peer_domain:     peer_domain,
		peer_localmspid: peer_localmspid,
		cryptopath:      cryptopath,
		configtxpath:    configtxpath,
		chaincodepath:   chaincodepath,
	}

	return tools
}
