package channel

import (
	"errors"
)

type ChannelManager struct {
	configPath   string
	channelinfos map[string]*Channelinfo
}

func NewChannelmanager(configPath string) *ChannelManager {
	cm := ChannelManager{}
	cm.configPath = configPath
	cm.channelinfos = make(map[string]*Channelinfo)
	return &cm
}

func (m *ChannelManager) NewChannelInfo(channelName string) (*Channelinfo, error) {
	if _, ok := m.channelinfos[channelName]; ok == true {
		return nil, errors.New("ChannelInfo already exists")
	}

	ci := Channelinfo{
		channelName: channelName,
		created:     false,
		configPath:  m.configPath,
	}

	m.channelinfos[channelName] = &ci
	return &ci, nil
}

func (m *ChannelManager) GetChannelInfo(channelname string) (*Channelinfo, bool) {
	ci, ok := m.channelinfos[channelname]
	return ci, ok
}

func (m *ChannelManager) GetChannelinfos() []*Channelinfo {
	var infos []*Channelinfo
	for _, ci := range m.channelinfos {
		infos = append(infos, ci)
	}
	return infos
}

func (m *ChannelManager) DeleteChannelInfo(channelname string) {
	delete(m.channelinfos, channelname)
}
