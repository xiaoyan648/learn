package errors

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

type WithCode struct {
	err   error
	code  int
	cause error
}

// Error return the externally-safe error message.
func (w *WithCode) Error() string { return fmt.Sprintf("%v", w.err) }

// Cause return the cause of the WithCode error.
func (w *WithCode) Cause() error { return w.cause }

// Unwrap provides compatibility for Go 1.13 error chains.
func (w *WithCode) Unwrap() error { return w.cause }

func WrapC(err error, code int) *WithCode {
	if err == nil {
		return nil
	}

	return &WithCode{
		err:   err,
		code:  code,
		cause: err,
	}
}

func WrapCE(err error, code int) error {
	if err == nil {
		return nil
	}

	return &WithCode{
		err:   err,
		code:  code,
		cause: err,
	}
}

func TestErrCompare(t *testing.T) {
	var err error
	we := WrapC(err, 1001)
	var e *WithCode

	// error 接口 和 *WithCode 的nil不想等
	assert.Equal(t, we, e)
	assert.Nil(t, we)
	assert.Nil(t, e)
	assert.Equal(t, err, nil)
	assert.NotEqual(t, err, we)
	var e2 error = we
	assert.NotEqual(t, err, e2)
	if we == nil {
		fmt.Println("aaa")
	}
	if e == nil {
		fmt.Println("aaa")
	}
	if e2 == nil { // 类型转化后判断 nil 不成功
		fmt.Println("aaa")
	}

	// error 接口 nil 与 nil 比较
	assert.Equal(t, err, WrapCE(err, 1001))
	assert.NoError(t, WrapCE(err, 1001))

}
