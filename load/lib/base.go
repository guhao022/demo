package lib

import "time"

// 原生请求的结构
type RawReq struct {
	Id  int64
	Req []byte
}

//原生相应的结构
type RawResp struct {
	Id     int64
	Resp   []byte
	Err    error
	Elapse time.Duration
}

// 相应代码
type ResultCode int

const (
	RESULT_SUCCESS              ResultCode = 0    //成功
	RESULT_WARNING_CALL_TIMEOUT ResultCode = 1001 //超时警告
	RESULT_ERROR_CALL           ResultCode = 2001 //调用错误
	RESULT_ERROR_RESPONSE       ResultCode = 2002 // 相应内容错误
	RESULT_ERROR_CALEE          ResultCode = 2003 // 被调用方（被测软件）的内部错误。
	RESULT_FATAL_CALL           ResultCode = 3001 // 调用过程中发生了致命错误！
)

func GetResultCodePlain(code ResultCode) string {
	var codePlain string
	switch code {
	case RESULT_SUCCESS:
		codePlain = "Success"
	case RESULT_WARNING_CALL_TIMEOUT:
		codePlain = "Call Timeout Warning"
	case RESULT_ERROR_CALL:
		codePlain = "Call Error"
	case RESULT_ERROR_RESPONSE:
		codePlain = "Response Error"
	case RESULT_ERROR_CALEE:
		codePlain = "Callee Error"
	case RESULT_FATAL_CALL:
		codePlain = "Call Fatal Error"
	default:
		codePlain = "Unknown result code"
	}
	return codePlain
}

type CallResult struct {
	Id     int64         //id
	Req    RawReq        //原生请求
	Resp   RawResp       //原生相应
	Code   ResultCode    //相应代码
	Msg    string        //简介
	Elapse time.Duration //耗时
}

// 载荷发生器的状态的类型
type GenStatus int

const (
	STATUS_ORIGINAL GenStatus = 0  //原始状态
	STATUS_STARTED GenStatus = 1	// 开启状态
	STATUS_STOPPED GenStatus = 2	// 停止状态
)

// 载荷发生器的接口。
type Generator interface {
	// 启动载荷发生器。
	Start()
	// 停止载荷发生器。
	// 第一个结果值代表已发载荷总数，且仅在第二个结果值为true时有效。
	// 第二个结果值代表是否成功将载荷发生器转变为已停止状态。
	Stop() (uint64, bool)
	// 获取状态。
	Status() GenStatus
}
