package webserver

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

var isAdmin bool = true

func mainPage(c *gin.Context) {
	if isAdmin {
		c.HTML(http.StatusOK, "IndexAdmin", nil)
	} else {
		c.HTML(http.StatusOK, "IndexUser", nil)
	}
}

func blockchainControl(c *gin.Context) {
	if isAdmin {
		c.HTML(http.StatusOK, "BlockchainControlAdmin", nil)
	} else {
		c.HTML(http.StatusOK, "BlockchainControlUser", nil)
	}
}

func logControl(c *gin.Context) {
	if isAdmin {
		c.HTML(http.StatusOK, "LogControlAdmin", nil)
	} else {
		c.HTML(http.StatusOK, "LogControlUser", nil)
	}
}

func networkControl(c *gin.Context) {
	if isAdmin {
		c.HTML(http.StatusOK, "NetworkControlAdmin", nil)
	} else {
		c.HTML(http.StatusOK, "NetworkControlUser", nil)
	}
}

func channelControl(c *gin.Context) {
	if isAdmin {
		c.HTML(http.StatusOK, "ChannelControlAdmin", nil)
	} else {
		c.HTML(http.StatusOK, "ChannelControlUser", nil)
	}
}

func chaincodeControl(c *gin.Context) {
	if isAdmin {
		c.HTML(http.StatusOK, "ChaincodeControlAdmin", nil)
	} else {
		c.HTML(http.StatusOK, "ChaincodeControlUser", nil)

	}
}

func userControl(c *gin.Context) {
	if isAdmin {
		c.HTML(http.StatusOK, "UserControlAdmin", nil)
	}
}

func log(c *gin.Context) {
	if isAdmin {
		c.HTML(http.StatusOK, "LogAdmin", nil)
	} else {
		c.HTML(http.StatusOK, "LogUser", nil)
	}
}

func login(c *gin.Context) {
	c.HTML(http.StatusOK, "Login.html", nil)
}

func userProfile(c *gin.Context) {
	c.HTML(http.StatusOK, "UserProfile.html", nil)
}

func message(c *gin.Context) {
	c.HTML(http.StatusOK, "Message.html", nil)
}

func addNetwork(c *gin.Context) {
	c.HTML(http.StatusOK, "addNetwork.html", nil)
}

func addBlockchain(c *gin.Context) {
	c.HTML(http.StatusOK, "addBlockchain.html", nil)
}

func addOrderer(c *gin.Context) {
	c.HTML(http.StatusOK, "addOrderer.html", nil)
}

func addPeer(c *gin.Context)  {
	c.HTML(http.StatusOK, "addPeer.html", nil)
}

func addChannel(c *gin.Context) {
	c.HTML(http.StatusOK, "addChannel.html", nil)
}

func updateChannel(c *gin.Context) {
	c.HTML(http.StatusOK, "updateChannel.html", nil)
}

func addContract(c *gin.Context) {
	c.HTML(http.StatusOK, "addContract.html", nil)
}

func invokeContract(c *gin.Context) {
	c.HTML(http.StatusOK, "invokeContract.html", nil)
}

func modifyContract(c *gin.Context) {
	c.HTML(http.StatusOK, "modifyContract.html", nil)
}

func instantiateContract(c *gin.Context) {
	c.HTML(http.StatusOK, "instantiateContract.html", nil)
}

func contractLog(c *gin.Context) {
	c.HTML(http.StatusOK, "contractLog.html", nil)
}
