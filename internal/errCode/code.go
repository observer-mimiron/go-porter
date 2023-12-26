package errCode

import (
	"fmt"
)

type ErrCode struct {
	hCode   int32  // http状态码
	bCode   int32  // 业务码
	message string // 错误描述
	alert   bool   // 警报通知
	error   error  // 错误
}

func NewErrCode(bCode int32, message string) *ErrCode {
	return &ErrCode{
		bCode:   bCode,
		message: message,
	}
}

func (e *ErrCode) Error() string {
	return fmt.Sprintf("ErrCode:%v, message:%s", e.hCode, e.message)
	//return errors.Wrapf(e.ErrCode), e.message).Error()
}
func (e *ErrCode) ErrCode() int32 {
	return e.bCode
}

func (e *ErrCode) Message() string {
	return e.message
}

func (e *ErrCode) Wrap(err error) *ErrCode {
	e.error = err
	return e
}

func (e *ErrCode) Alert() *ErrCode {
	e.alert = true
	return e
}
