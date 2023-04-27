package channel

import (
	"bit-bass/deploy"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
)

func WriteSdkConfig(configurator *deploy.Configurator, cm *ChannelManager) error {
	//SdkConfig文件中的基本内容
	config := map[string]interface{}{
		"version": "1.0.0",
		"client": map[string]interface{}{
			"logging": map[string]interface{}{
				"level": "error",
			},
			"cryptoconfig": map[string]interface{}{
				"path": "../artifacts/crypto-info",
			},
			"BCCSP": map[string]interface{}{
				"security": map[string]interface{}{
					"enabled": true,
					"default": map[string]interface{}{
						"provider": "SW",
					},
					"hashAlgorithm": "SHA2",
					"softVerify":    true,
					"level":         256,
				},
			},
		},
	}
	//SdkConfig中需要配置的部分
	organizations := make(map[string]interface{})
	orderers := make(map[string]interface{})
	peers := make(map[string]interface{})
	channels := make(map[string]interface{})

	//将peer组织中的信息加入到配置中
	for _, peerorg := range configurator.Peerorgsconf {
		var peerdomins []string
		for _, peer := range peerorg.GetNodes() {
			peerdomin := peer.Host() + "." + peer.Domain().Domain()
			peerdomins = append(peerdomins, peerdomin)
			peers[peerdomin] = map[string]interface{}{
				"url": "localhost:" + peer.Port().Port(),
				"grpcOptions": map[string]interface{}{
					"allow-insecure": true,
				},
			}
		}

		organizations[peerorg.Name()] = map[string]interface{}{
			"mspid": peerorg.MSPID(),
			"cryptoPath": peerorg.CryptoPath() + string(os.PathSeparator) + "users" + string(os.PathSeparator) +
				"Admin@" + peerorg.Domain() + string(os.PathSeparator) + "msp",
			"peers": peerdomins,
		}
	}
	//将orderer组织中到信息加入到配置
	organizations[configurator.Ordererorgconf.Name()] = map[string]interface{}{
		"mspID": configurator.Ordererorgconf.MSPID(),
		"cryptoPath": configurator.Ordererorgconf.CryptoPath() + string(os.PathSeparator) +
			"users" + string(os.PathSeparator) + "Admin@" + configurator.Ordererorgconf.Domain() +
			string(os.PathSeparator) + "msp",
	}
	ordererorg := configurator.Ordererorgconf
	orderers[ordererorg.Name()+"."+ordererorg.Domain()] = map[string]interface{}{
		"url": "localhost:" + configurator.Ordererorgconf.GetNodes()[0].Port().Port(),
		"grpcOptions": map[string]interface{}{
			"allow-insecure": true,
		},
	}

	channelinfos := cm.GetChannelinfos()
	for _, channelinfo := range channelinfos {
		channelpeers := make(map[string]interface{})
		for _, peer := range channelinfo.GetPeers() {
			channelpeers[peer] = map[string]interface{}{
				"endorsingPeer":  true,
				"chaincodeQuery": true,
				"ledgerQuery":    true,
				"eventSource":    true,
			}
		}
		channels[channelinfo.GetChannelName()] = map[string]interface{}{
			"peers":    channelpeers,
			"orderers": []string{configurator.Ordererorgconf.Name() + "." + configurator.Ordererorgconf.Domain()},
		}
	}
	config["channels"] = channels
	config["organizations"] = organizations
	config["orderers"] = orderers
	config["peers"] = peers

	configbuffer, err := yaml.Marshal(config)
	if err != nil {
		return errors.Errorf("Failed to marshal the info %v :", err)
	}

	err = ioutil.WriteFile(configurator.Outpath+"/sdk-config.yaml", configbuffer, 0666)
	if err != nil {
		return errors.Errorf("Faile to write info to file %v :", err)
	}

	return nil
}
