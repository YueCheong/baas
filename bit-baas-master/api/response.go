package api

//后端对请求返回的回执
type Response struct {
	//回执的错误代码，200表示正常
	Code ApiResponseCode
	//回执的错误信息，
	Msg     string
	Package interface{}
}

const (
	SUCCESS         ApiResponseCode = 200
	PARTIAL_CONTENT ApiResponseCode = 206
	ERROR           ApiResponseCode = 500
	INVALID_PARAMS  ApiResponseCode = 400
	NOT_IMPLEMENTED ApiResponseCode = 501

	INVALID_OPERATION         ApiResponseCode = 10001
	BLOCKCHAIN_NETWORK_ERROR  ApiResponseCode = 10002
	BLOCKCHAIN_CHANNEL_ERROR  ApiResponseCode = 10003
	BLOCKCHAIN_CONTRACT_ERROR ApiResponseCode = 10004
	RESOURCES_NOT_FOUND       ApiResponseCode = 10005
	DOCKER_NETWORK_ERROR      ApiResponseCode = 10006
)

var apiResponseMsg = map[ApiResponseCode]string{
	SUCCESS:         "ok",
	PARTIAL_CONTENT: "部分资源请求成功",
	ERROR:           "error",
	INVALID_PARAMS:  "请求参数错误",
	NOT_IMPLEMENTED: "功能未实现",

	INVALID_OPERATION:         "非法操作",
	BLOCKCHAIN_NETWORK_ERROR:  "区块链网络错误",
	BLOCKCHAIN_CHANNEL_ERROR:  "区块链通道错误",
	BLOCKCHAIN_CONTRACT_ERROR: "区块链链码错误",
	RESOURCES_NOT_FOUND:       "找不到请求的资源",
	DOCKER_NETWORK_ERROR:      "Docker网络错误",
}

type ApiResponseCode int

//获取返回错误码对应的信息
func (code ApiResponseCode) GetMsg() string {
	return GetApiResponseMsg(code)
}

func (code ApiResponseCode) Int() int {
	return int(code)
}

//获取返回错误码对应的信息
func GetApiResponseMsg(code ApiResponseCode) string {
	msg, ok := apiResponseMsg[code]
	if ok {
		return msg
	}
	return "未知错误代码"
}
