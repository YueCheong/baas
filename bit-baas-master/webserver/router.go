package webserver

import (
	"bit-bass/api"
	"github.com/gin-contrib/multitemplate"
	"github.com/gin-gonic/gin"
	"path/filepath"
	"strings"
)

func InitRouter() *gin.Engine {
	r := gin.New()

	r.Use(gin.Logger())

	r.Use(gin.Recovery())

	//读取模板
	r.HTMLRender = loadTemplates("./webserver/templates")
	//设置静态资源
	r.Static("static/", "./webserver/static")
	//读取网页

	//注册网页路由
	dashboard := r.Group("")
	{
		dashboard.GET("/", mainPage)
		dashboard.GET("/network", networkControl)
		dashboard.GET("/blockchain", blockchainControl)
		dashboard.GET("/channel", channelControl)
		dashboard.GET("/chaincode", chaincodeControl)
		dashboard.GET("/user", userControl)
		dashboard.GET("/login", login)
		dashboard.GET("/log", logControl)
		dashboard.GET("/profile", userProfile)
		dashboard.GET("/message", message)
	}

	//注册api路由
	apiGroup := r.Group("/api")
	{
		apiGroup.GET("/summary", api.GetSummary)

		apiGroup.GET("/networks", api.GetNetworks)
		apiGroup.PUT("/networks/:name", api.CreatNetwork)
		apiGroup.DELETE("/networks/:name", api.RemoveNetwork)

		apiGroup.GET("/blockchains", api.GetBlockchains)
		apiGroup.PUT("/blockchains", api.CreatBlockchain)
		apiGroup.POST("/blockchains", api.ManageBlockchain)
		apiGroup.DELETE("/blockchains/:id", api.DeleteBlockchain)

		apiGroup.GET("/channels", api.GetChannels)
		apiGroup.PUT("/channels", api.CreatChannel)
		apiGroup.POST("/channels", api.ManageChannel)

		apiGroup.GET("/contracts", api.GetContracts)
		apiGroup.POST("/newcontract", api.CreatContract)
		apiGroup.POST("/contracts", api.ManageContract)
		apiGroup.POST("/contractcall", api.InvokeContract)
		apiGroup.GET("/contractlog/:id", api.GetContractLogs)
		apiGroup.GET("/contractlog", api.GetAllContractLogs)
	}

	// pwindow
	popWindow := r.Group("/pop")
	{
		popWindow.GET("/addNetwork", addNetwork)
		popWindow.GET("/addBlockchain", addBlockchain)
		popWindow.GET("/addOrderer", addOrderer)
		popWindow.GET("/addPeer", addPeer)
		popWindow.GET("/addChannel", addChannel)
		popWindow.GET("/updateChannel", updateChannel)
		popWindow.GET("/addContract", addContract)
		popWindow.GET("/invokeContract", invokeContract)
		popWindow.GET("/modifyContract", modifyContract)
		popWindow.GET("/instantiateContract", instantiateContract)
		popWindow.GET("/contractLog", contractLog)
	}
	return r

}

func loadTemplates(templatesDir string) multitemplate.Renderer {
	r := multitemplate.NewRenderer()

	//读入子页面
	subPages, err := filepath.Glob(templatesDir + "/subPages/*.html")
	if err != nil {
		panic(err.Error())
	}

	//读入User管理页面框架
	userFrame, err := filepath.Glob(templatesDir + "/frame/MainFrameUser.html")
	if err != nil {
		panic(err.Error())
	}
	//注册继承了User页面框架的页面
	for _, subPage := range subPages {
		layoutCopy := make([]string, len(userFrame))
		copy(layoutCopy, userFrame)
		files := append(layoutCopy, subPage)
		pageName := strings.TrimSuffix(filepath.Base(subPage), filepath.Ext(subPage))
		r.AddFromFiles(pageName+"User", files...)
	}
	//读入Admin管理页面框架
	adminFrame, err := filepath.Glob(templatesDir + "/frame/MainFrameAdmin.html")
	if err != nil {
		panic(err.Error())
	}

	//注册继承了Admin页面框架的页面
	for _, subPage := range subPages {
		layoutCopy := make([]string, len(adminFrame))
		copy(layoutCopy, adminFrame)
		tempFiles := append(layoutCopy, subPage)
		pageName := strings.TrimSuffix(filepath.Base(subPage), filepath.Ext(subPage))
		r.AddFromFiles(pageName+"Admin", tempFiles...)
	}

	//读入空白页面框架
	emptyFrame, err := filepath.Glob(templatesDir + "/frame/EmptyBase.html")
	if err != nil {
		panic(err.Error())
	}
	//读入独立页面
	Pages, err := filepath.Glob(templatesDir + "/pages/*.html")
	if err != nil {
		panic(err.Error())
	}
	//注册独立页面
	for _, Page := range Pages {
		layoutCopy := make([]string, len(emptyFrame))
		copy(layoutCopy, emptyFrame)
		files := append(layoutCopy, Page)
		r.AddFromFiles(filepath.Base(Page), files...)
	}

	layouts, err := filepath.Glob(templatesDir + "/layouts/empty.html")
	if err != nil {
		panic(err.Error())
	}
	includes, err := filepath.Glob(templatesDir + "/popPages/*.html")
	if err != nil {
		panic(err.Error())
	}
	// 为layouts/和includes/目录生成 templates map
	for _, include := range includes {
		layoutCopy := make([]string, len(layouts))
		copy(layoutCopy, layouts)
		files := append(layoutCopy, include)
		r.AddFromFiles(filepath.Base(include), files...)
	}
	return r
}
