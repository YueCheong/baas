package artifacts

import (
	"bit-bass/deploy"
	"bytes"
	"fmt"
	"github.com/pkg/errors"
	"io/ioutil"
	"os"
	"os/exec"
	"os/user"
	"runtime"
	"strconv"
)

//用于生成配置及批处理文件
func GenerateYaml(configurator *deploy.Configurator, path string) error {
	err := os.MkdirAll(path, 0777)
	if err != nil {
		return err
	}
	if err := GenerateCryptoYaml(configurator, path); err != nil {
		return err
	}
	if err := GenerateConfigtxYaml(configurator, path); err != nil {
		return err
	}
	return nil
}

//用于生成crypto-config
func GenerateCryptoYaml(configurator *deploy.Configurator, path string) error {
	oconfig := configurator.Ordererorgconf
	var buffer bytes.Buffer
	buffer.WriteString("OrdererOrgs:\n")
	buffer.WriteString("  - Name: ")
	buffer.WriteString(oconfig.Name())
	buffer.WriteString("\n    Domain: ")
	buffer.WriteString(oconfig.Domain())
	buffer.WriteString("\n    Specs:\n")
	buffer.WriteString("      - Hostname: ")
	buffer.WriteString(oconfig.GetNodes()[0].Host())
	buffer.WriteString("\nPeerOrgs:\n")
	for _, pconfig := range configurator.Peerorgsconf {
		buffer.WriteString("  - Name: ")
		buffer.WriteString(pconfig.Name())
		buffer.WriteString("\n    Domain: ")
		buffer.WriteString(pconfig.Domain())
		buffer.WriteString("\n    Template: ")
		buffer.WriteString("\n      Count: ")
		buffer.WriteString(strconv.Itoa(len(pconfig.GetNodes())))
		buffer.WriteString("\n    Users: ")
		buffer.WriteString("\n      Count: ")
		buffer.WriteString("1")
		buffer.WriteString("\n")
	}
	//print(buffer.String())
	err := ioutil.WriteFile(path+"/crypto-config.yaml", buffer.Bytes(), 0777)
	return err
}

//用于生成configtx
func GenerateConfigtxYaml(configurator *deploy.Configurator, path string) error {
	oconfig := configurator.Ordererorgconf
	var buffer bytes.Buffer
	buffer.WriteString("Organizations:")
	buffer.WriteString("\n    - &")
	buffer.WriteString(oconfig.Name())
	buffer.WriteString("Org")
	buffer.WriteString("\n        Name: ")
	buffer.WriteString(oconfig.Name())
	buffer.WriteString("\n        ID: ")
	buffer.WriteString(oconfig.MSPID())
	buffer.WriteString("\n        MSPDir: ")
	buffer.WriteString("crypto-config/ordererOrganizations/")
	buffer.WriteString(oconfig.Domain())
	buffer.WriteString("/msp")
	for _, pconfig := range configurator.Peerorgsconf {
		anchor := pconfig.GetAnchorPeer()
		if anchor == nil {
			anchor = pconfig.GetNodes()[0]
		} else {
			fmt.Println(anchor)
			fmt.Println("Set anchor peer :", anchor.Host()+pconfig.Domain())
		}

		buffer.WriteString("\n    - &")
		buffer.WriteString(pconfig.Name())
		buffer.WriteString("Org")
		buffer.WriteString("\n        Name: ")
		buffer.WriteString(pconfig.Name())
		buffer.WriteString("\n        ID: ")
		buffer.WriteString(pconfig.MSPID())
		buffer.WriteString("\n        MSPDir: ")
		buffer.WriteString("crypto-config/peerOrganizations/")
		buffer.WriteString(pconfig.Domain())
		buffer.WriteString("/msp")
		buffer.WriteString("\n        AnchorPeers: ")
		buffer.WriteString("\n          - Host: ")
		buffer.WriteString(anchor.Host() + ".")
		buffer.WriteString(pconfig.Domain())
		buffer.WriteString("\n            Port: ")
		buffer.WriteString(strconv.Itoa(anchor.Port().Int()))
	}

	buffer.WriteString("\nOrderer: &OrdererDefaults")
	buffer.WriteString("\n    Orderertype: solo")
	buffer.WriteString("\n    Addresses: ")
	buffer.WriteString("\n      - orderer.")
	buffer.WriteString(oconfig.Domain())
	buffer.WriteString(":")
	buffer.WriteString(strconv.Itoa(oconfig.GetNodes()[0].Port().Int()))
	buffer.WriteString("\n    BatchTimeout: 2s")
	buffer.WriteString("\n    BatchSize: ")
	buffer.WriteString("\n        MaxMessageCount: 10")
	buffer.WriteString("\n        AbsoluteMaxBytes: 99 MB")
	buffer.WriteString("\n        PreferredMaxBytes: 512 KB")

	buffer.WriteString("\nProfiles:")
	buffer.WriteString("\n    TwoOrgOrdererGenesis:")
	buffer.WriteString("\n        Orderer:")
	buffer.WriteString("\n            <<: *OrdererDefaults")
	buffer.WriteString("\n            Organizations:")
	buffer.WriteString("\n                - *")
	buffer.WriteString(oconfig.Name())
	buffer.WriteString("Org")
	buffer.WriteString("\n        Consortiums:")
	buffer.WriteString("\n            SampleConsortium:")
	buffer.WriteString("\n                Organizations:")
	for _, pconfig := range configurator.Peerorgsconf {
		buffer.WriteString("\n                    - *")
		buffer.WriteString(pconfig.Name())
		buffer.WriteString("Org")
	}
	buffer.WriteString("\n    TwoOrgChannel:")
	buffer.WriteString("\n        Consortium: SampleConsortium")
	buffer.WriteString("\n        Application:")
	buffer.WriteString("\n            Organizations:")
	for _, pconfig := range configurator.Peerorgsconf {
		buffer.WriteString("\n                    - *")
		buffer.WriteString(pconfig.Name())
		buffer.WriteString("Org")
	}
	buffer.WriteString("\n")
	//print(buffer.String())
	err := ioutil.WriteFile(path+"/configtx.yaml", buffer.Bytes(), 0777)
	return err
}

//用于生成批处理文件
//不用了，只作参考
//func WriteBat(configurator *deploy.Configurator, path string) error {
//	var buffer bytes.Buffer
//	buffer.WriteString("if not exist configtx (mkdir configtx)\n\n")
//	buffer.WriteString("docker-compose run tool cryptogen generate --config=/tmp/crypto-config.yaml\n\n")
//	buffer.WriteString("docker-compose run tool configtxgen -profile TwoOrgOrdererGenesis -outputBlock
//	/tmp/configtx/genesis.block\n\n")
//	buffer.WriteString("docker-compose run tool configtxgen -profile TwoOrgChannel -outputCreateChannelTx
//	/tmp/configtx/channel.tx -channelID mychannel\n\n")
//	for _, pconfig := range configurator.Peerorgsconf {
//		buffer.WriteString("docker-compose run tool configtxgen -profile TwoOrgChannel
//		-outputAnchorPeersUpdate /tmp/configtx/" +
//			pconfig.Name() +
//			"MSPanchors.tx -channelID mychannel -asOrg ")
//		buffer.WriteString(pconfig.Name())
//		buffer.WriteString("\n\n")
//	}
//	//print(buffer.String())
//	err := ioutil.WriteFile(path + "/run.bat", buffer.Bytes(), 0777)
//	return err
//}

//用于执行批处理文件
//不用了，只作参考
//func ExecuteBat(path string) error {
//	_ = os.RemoveAll(path + string(os.PathSeparator) + "crypto-config")
//	commandName := path + string(os.PathSeparator) + "run.bat"
//	cmd := exec.Command(commandName)
//	cmd.Stderr = os.Stderr
//	cmd.Dir = path
//	//fmt.Println("cmd", cmd.Path, cmd.Dir, cmd.Args)
//	err := cmd.Run()
//	if err != nil {
//		return err
//	}
//	err = cmd.Wait()
//	return nil
//}

func GenerateCryptoConfig(path string) error {
	_ = os.RemoveAll(path + string(os.PathSeparator) + "crypto-config")

	//Docker默认以root身份运行container，会导致通过tool镜像创建的crypto文件及
	//tx文件属于root用户，导致权限问题，因此需要当前用户的用户id来以用户身份启动
	//tool镜像

	var args []string
	args = append(args, "run", "-v", path+"/:/etc/bit-baas/artifacts/")

	if runtime.GOOS == "linux" {
		userId, err := user.Current()
		if err != nil {
			return errors.Errorf("Failed to generate crypto-config"+
				": Can't get user info : %v ", err)
		}
		args = append(args, "-u", userId.Uid+":"+userId.Uid)
	}

	args = append(args, "tool", "cryptogen", "generate",
		"--config=/etc/bit-baas/artifacts/crypto-config.yaml",
		"--output=/etc/bit-baas/artifacts/crypto-config")

	//命令说明：
	//通过tool镜像中的cryptogen工具，通过crypto-config.yaml文件生成区块链系统所需的证书文件
	//启动tool镜像中的选项-v表示将宿主机中的路径挂在到容器中，选项-u表示指定容器运行时使用的用户
	//如果被指定用户，容器将默认以root身份运行，导致创建的证书文件属于root用户，产生权限问题
	cryptoGenCMD := exec.Command("docker-compose", args...)

	var stdout, stderr bytes.Buffer
	cryptoGenCMD.Stdout = &stdout
	cryptoGenCMD.Stderr = &stderr

	err := cryptoGenCMD.Run()
	//fmt.Printf("generate crypto-config stdout:%v\n",stdout.String())
	//fmt.Printf("generate crypto-config stderr:%v\n",stderr.String())
	if err != nil {
		return errors.Errorf("Failed to generate crypto config : Failed to exec command : %v\n"+
			"The stderr out put is : \n %v \n", err, stderr.String())
	}

	return nil
}

//fabric联盟链启动时必须为每个Orderer节点提供一个创世块。Orderer将在这个创世块上启动系统通道
//我们创建的应用通道时提交创建交易后，系统会自动为其创建创世块，因此不需要我们自己创建
func GenerateGenesisBlock(path string) error {
	var args []string
	args = append(args, "run", "-v", path+"/:/etc/bit-baas/artifacts/")

	if runtime.GOOS == "linux" {
		userId, err := user.Current()
		if err != nil {
			return errors.Errorf("Failed to generate crypto-config"+
				": Can't get user info : %v ", err)
		}
		args = append(args, "-u", userId.Uid+":"+userId.Uid)
	}

	args = append(args, "tool", "configtxgen", "-configPath",
		"/etc/bit-baas/artifacts/", "-profile", "TwoOrgOrdererGenesis",
		"-channelID", "orderer-system-channel", "-outputBlock",
		"/etc/bit-baas/artifacts/configtx/genesis.block")

	err := os.MkdirAll(path+"/configtx", 0777)

	GenBlockCMD := exec.Command("docker-compose", args...)

	var stdout, stderr bytes.Buffer
	GenBlockCMD.Stdout = &stdout
	GenBlockCMD.Stderr = &stderr

	err = GenBlockCMD.Run()
	//fmt.Printf("generate GenesisBlock stdout:%v\n",stdout.String())
	//fmt.Printf("generate GenesisBlock stderr:%v\n",stderr.String())
	if err != nil {
		return errors.Errorf("Failed to generate genesis block : Failed to exec command : %v\n"+
			"The stderr out put is : \n %v \n", err, stderr.String())
	}

	return nil
}

//创建的通道建立交易文件位置为 “path/configtx/chName.tx” 其中path是传入的参数文件目录
//chName时传入的参数通道名
func GenerateChannelCreationTx(path string, chName string) (string, error) {
	var args []string
	args = append(args, "run", "-v", path+"/:/etc/bit-baas/artifacts/")

	if runtime.GOOS == "linux" {
		userId, err := user.Current()
		if err != nil {
			return "", errors.Errorf("Failed to generate crypto-config"+
				": Can't get user info : %v ", err)
		}
		args = append(args, "-u", userId.Uid+":"+userId.Uid)
	}

	//生成的文件相对于config路径的地址
	txRelDir := "/configtx/" + chName + ".tx"

	args = append(args, "tool", "configtxgen", "-configPath",
		"/etc/bit-baas/artifacts/", "-profile", "TwoOrgChannel", "-channelID",
		chName, "-outputCreateChannelTx", "/etc/bit-baas/artifacts"+txRelDir)

	err := os.MkdirAll(path+"/configtx", 0777)

	GenTxCMD := exec.Command("docker-compose", args...)

	var stdout, stderr bytes.Buffer
	GenTxCMD.Stdout = &stdout
	GenTxCMD.Stderr = &stderr

	err = GenTxCMD.Run()
	//fmt.Printf("generate ChannelCreationTx stdout: %v\n",stdout.String())
	//fmt.Printf("generate ChannelCreationTx stderr:%v \n",stderr.String())
	if err != nil {
		return "", errors.Errorf("Failed to generate channel creation tx : Failed to exec command : %v\n"+
			"The stderr out put is : \n %v \n", err, stderr.String())
	}

	return path + txRelDir, nil
}

func GenerateAnchorPeerTx(path string, chName, orgName string) (string, error) {
	var args []string
	args = append(args, "run", "-v", path+"/:/etc/bit-baas/artifacts/")

	if runtime.GOOS == "linux" {
		userId, err := user.Current()
		if err != nil {
			return "", errors.Errorf("Failed to generate crypto-config"+
				": Can't get user info : %v ", err)
		}
		args = append(args, "-u", userId.Uid+":"+userId.Uid)
	}

	//生成的文件相对于config路径的地址
	txRelDir := "/configtx/" + "anchor_" + chName + "_" + orgName + ".tx"

	args = append(args, "tool", "configtxgen", "-configPath", "/etc/bit-baas/artifacts/",
		"-profile", "TwoOrgChannel", "-channelID", chName, "-asOrg", orgName,
		"-outputAnchorPeersUpdate", "/etc/bit-baas/artifacts"+txRelDir)

	//如果路径生成路径不存在则创建路径
	err := os.MkdirAll(path+"/configtx", 0777)

	GenTxCMD := exec.Command("docker-compose", args...)

	var stdout, stderr bytes.Buffer
	GenTxCMD.Stdout = &stdout
	GenTxCMD.Stderr = &stderr

	err = GenTxCMD.Run()
	//fmt.Printf("generate ChannelCreationTx stdout:%v\n",stdout.String())
	//fmt.Printf("generate ChannelCreationTx stderr:%v\n",stderr.String())

	if err != nil {
		return "", errors.Errorf("Failed to generate anchor tx : Failed to exec command : %v\n"+
			"The stderr out put is : \n %v \n", err, stderr.String())
	}

	return path + txRelDir, nil
}
