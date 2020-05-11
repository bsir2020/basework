package api

type Message struct {
	Item    string        `json:"item"`
	Subject string        `json:"subject"`
	Mtype   int           `json:"type"` //#0发起 1反馈
	Status  bool          `json:"status"`
	Id      int64         `json:"id"`
	Data    []interface{} `json:"data"`
}

// 定义错误码
type Errno struct {
	Code    int
	Message string
}

//各系统模块错误码范围
//1000-1999 系统框架
//2000-2999 游戏大厅
//3000-3999 游戏
//4000-4999 基础后台
var (
	OK = &Errno{Code: 0, Message: "OK"}

	// 系统错误, 前缀为 100
	SysInternalServerError = &Errno{Code: 1001, Message: "内部服务器错误"}
	SysErrEncrypt          = &Errno{Code: 1004, Message: "加密用户密码时发生错误"}

	// 数据库错误, 前缀为 101
	DBErr     = &Errno{Code: 1011, Message: "数据库错误"}
	DBErrFill = &Errno{Code: 1012, Message: "从数据库填充 struct 时发生错误"}

	// 认证错误, 前缀是 102
	AuthErr = &Errno{Code: 1021, Message: "验证失败"}
	AuthExp = &Errno{Code: 1022, Message: "jwt 是无效的"}

	//redis错误
	RedisErr = &Errno{Code: 1011, Message: "redis connection error"}
	RedisRun = &Errno{Code: 1012, Message: "从数据库填充 struct 时发生错误"}

	//mq错误
	MQConnErr = &Errno{Code: 1011, Message: "数据库错误"}
	MQRun     = &Errno{Code: 1012, Message: "从数据库填充 struct 时发生错误"}
)
