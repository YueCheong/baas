package network

import (
	"context"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

type NodesNet struct {
	ctx  context.Context
	cli  *client.Client
	id   string
	name string
}

func (net *NodesNet) NetName() string {
	return net.name
}

func (net *NodesNet) NetID() string {
	return net.id
}

func (net *NodesNet) CreateNet() error {
	resp, err := net.cli.NetworkCreate(net.ctx, net.name, types.NetworkCreate{})
	if err != nil {
		return err
	}
	net.id = resp.ID
	return nil
}

func (net *NodesNet) RemoveNet() error {
	err := net.cli.NetworkRemove(net.ctx, net.id)
	if err != nil {
		return err
	}
	return nil
}

func (net *NodesNet) IsNetExist() bool {
	_, err := net.cli.NetworkInspect(net.ctx, net.id)
	if err != nil {
		return false
	}
	return true
}

func (net *NodesNet) Inspect() types.NetworkResource {
	resp, err := net.cli.NetworkInspect(net.ctx, net.id)
	if err != nil {
		panic(err)
	}
	return resp
}

func NewNodesNet(netname string) *NodesNet {
	ctx := context.Background()
	cli, err := client.NewEnvClient()
	if err != nil {
		panic(err)
	}
	return &NodesNet{
		ctx:  ctx,
		cli:  cli,
		name: netname,
	}
}

func NewDefaultNodesNet(netname string) *NodesNet {
	ctx := context.Background()
	cli, err := client.NewEnvClient()
	if err != nil {
		panic(err)
	}
	return &NodesNet{
		ctx: ctx,
		cli: cli,
		//用标识mark来区分不同的网络
		name: netname,
	}
}
