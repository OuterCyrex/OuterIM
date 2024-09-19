package xerr

var codeText = map[int]string{
	SERVER_COMMON_ERROR: "服务器异常",
	REQUEST_PARAM_ERROR: "请求参数有误",
	DB_ERROR:            "数据库正忙",
}

func ErrMsg(errCode int) string {
	if msg, ok := codeText[errCode]; ok {
		return msg
	}
	return "unknown error"
}
