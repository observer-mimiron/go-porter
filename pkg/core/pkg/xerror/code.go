package xerror

import (
	"fmt"
)

type ErrCode struct {
	errCode int    // 业务码
	errMsg  string // 错误描述
}

func Error(errCode int, errMsg string) *ErrCode {
	return &ErrCode{
		errCode: errCode,
		errMsg:  errMsg,
	}
}

func (e *ErrCode) Error() string {
	return fmt.Sprintf("ErrCode:%d, ErrMsg:%s", e.errCode, e.errMsg)
	//return errors.Wrapf(e.ErrCode), e.message).Error()
}
func (e *ErrCode) Code() int {
	return e.errCode
}

func (e *ErrCode) Message() string {
	return e.errMsg
}
