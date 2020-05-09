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

//1000-2000 系统框架
//2000-3000 游戏大厅
//4000-5000 游戏
//5000-6000 基础后台
var (
	OK = &Errno{Code: 0, Message: "OK"}

	// 系统错误, 前缀为 100
	InternalServerError = &Errno{Code: 1001, Message: "内部服务器错误"}
	ErrBind             = &Errno{Code: 1002, Message: "请求参数错误"}
	ErrTokenSign        = &Errno{Code: 1003, Message: "签名 jwt 时发生错误"}
	ErrEncrypt          = &Errno{Code: 1004, Message: "加密用户密码时发生错误"}

	// 数据库错误, 前缀为 101
	ErrDatabase = &Errno{Code: 1011, Message: "数据库错误"}
	ErrFill     = &Errno{Code: 1012, Message: "从数据库填充 struct 时发生错误"}

	// 认证错误, 前缀是 102
	ErrValidation   = &Errno{Code: 1021, Message: "验证失败"}
	ErrTokenInvalid = &Errno{Code: 1022, Message: "jwt 是无效的"}

	//redis错误
	RedisErr = &Errno{Code: 1011, Message: "数据库错误"}
	RedisRun = &Errno{Code: 1012, Message: "从数据库填充 struct 时发生错误"}

	//mq错误
	MQConnErr = &Errno{Code: 1011, Message: "数据库错误"}
	MQRun     = &Errno{Code: 1012, Message: "从数据库填充 struct 时发生错误"}
)
