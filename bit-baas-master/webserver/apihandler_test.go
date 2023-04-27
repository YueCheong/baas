package webserver

import (
	"bit-bass/api"
	"bit-bass/core"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestApi(t *testing.T) {
	core.InitCore()
	r := InitRouter()

	testCreatNetwork(r, t)
	testGetNetworks(r, t)

	//关闭服务器前关闭所有网络
	for _, b := range core.CoreBlockChainManager.GetBlockchains() {
		core.CoreBlockChainManager.StopAndRemoveBlockchainById(b.GetId())
	}

}

func testCreatNetwork(r *gin.Engine, t *testing.T) {
	net := core.BlockchainToJson{
		Name: "firstnet",
		OrdererOrg: core.OrdererOrgToJson{
			Name:   "Orderer",
			Domain: "Testing.com",
			MSPID:  "OrdererMSP",
			Orderer: core.NodeToJson{
				Name: "Orderer",
				Port: "5050",
			},
		},
		PeerOrg: []core.PeerOrgToJson{{
			Name:   "org1",
			Domain: "Testing.com",
			MSPID:  "Org1MSP",
			Peers: []core.NodeToJson{
				{
					Name: "peer0",
					Port: "5051",
				},
				{
					Name: "peer1",
					Port: "6051",
				},
			},
		}, {
			Name:   "org2",
			Domain: "Testing.com",
			MSPID:  "Org2MSP",
			Peers: []core.NodeToJson{
				{
					Name: "peer0",
					Port: "7051",
				}, {
					Name: "peer1",
					Port: "8051",
				},
			},
		},
		},
		Config: nil,
		Channels: []core.ChannelToJson{
			{
				Name: "mych",
				Peers: []string{"peer0.org1.Testing.com", "peer1.org1.Testing.com",
					"peer0.org2.Testing.com", "peer1.org2.Testing.com"},
			},
		},
	}

	reqParm, err := json.Marshal(net)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(reqParm))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/api/networks", bytes.NewBuffer(reqParm))

	r.ServeHTTP(w, req)

	if w.Code != api.SUCCESS.Int() {
		t.Fail()
	}
	fmt.Println(w)
}

func testGetNetworks(r *gin.Engine, t *testing.T) {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/networks", nil)

	r.ServeHTTP(w, req)

	if w.Code != api.SUCCESS.Int() {
		t.Fail()
	}

	fmt.Println(w)
}
