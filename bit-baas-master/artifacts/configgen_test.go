package artifacts

import (
	"bit-bass/deploy"
	"bit-bass/utils"
	"fmt"
	"testing"
)

const id = 260

func TestGenerateConfig(t *testing.T) {

	configurator := *deploy.NewConfigurator(utils.ConfigPath())
	configurator.SetOrdererOrg("Orderer", "example.com", "OrdererMSP")
	configurator.SetOrdererNode("orderer", "7050/tcp")
	configurator.AddPeerOrg("Org1", "org1.example.com", "Org1MSP")
	configurator.AddPeerOrg("Org2", "org2.example.com", "Org2MSP")
	configurator.AddPeerNodeToOrg("Org1", "peer0", "7051/tcp")
	configurator.AddPeerNodeToOrg("Org1", "peer1", "8051/tcp")
	configurator.AddPeerNodeToOrg("Org2", "peer0", "9051/tcp")
	configurator.AddPeerNodeToOrg("Org2", "peer1", "10051/tcp")
	configurator.SetAnchorPeerToOrg("Org1", "peer1")
	configurator.SetAnchorPeerToOrg("Org2", "peer1")

	GenerateYaml(&configurator, utils.ConfigPathWithId(id))

	err := GenerateCryptoConfig(utils.ConfigPathWithId(id))
	if err != nil {
		t.Fatal(err)
	}

	err = GenerateGenesisBlock(utils.ConfigPathWithId(id))
	if err != nil {
		t.Fatal(err)
	}

	txdir, err := GenerateChannelCreationTx(utils.ConfigPathWithId(id), "mych")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("the tx dir is : ", txdir)

	txdir, err = GenerateAnchorPeerTx(utils.ConfigPathWithId(id), "mych", "Org1")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("the anchor dir is : ", txdir)
}
