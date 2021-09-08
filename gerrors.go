package gerrors

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"runtime"
	"time"
)

type CustomError struct {
	cause string
	msg   map[string]interface{}
	stack stack
}

type stack struct {
	funcName string
	file     string
}

func New(cause string) error {
	return &CustomError{
		cause: cause,
		msg:   make(map[string]interface{}),
		stack: caller(),
	}
}

func (c *CustomError) Error() string {
	return c.cause
}
func Trace(e error) {
	usecaseName := usecaseCaller()
	var ce *CustomError
	enc := json.NewEncoder(os.Stdout)
	if errors.As(e, &ce) {
		format := map[string]interface{}{
			"usecase": usecaseName,
			"func":    ce.stack.funcName,
			"file":    ce.stack.file,
			"msg":     ce.cause,
			"time":    time.Now(),
			"level":   "error",
		}
		enc.Encode(format)
	} else {
		format := map[string]interface{}{
			"usecase": usecaseName,
			"msg":     e.Error(),
			"time":    time.Now(),
			"level":   "error",
		}
		enc.Encode(format)
	}
}

func Wrap(e error, value map[string]interface{}) error {
	var ce *CustomError
	if errors.As(e, &ce) {
		for k, v := range value {
			ce.msg[k] = v
		}
		return ce
	}
	return e
}

func UnWrap(e error, key string) interface{} {
	var ce *CustomError
	if errors.As(e, &ce) {
		if v, ok := ce.msg[key]; ok {
			return v
		}
	}
	return nil
}

func caller() stack {
	pc, file, line, _ := runtime.Caller(2)
	funcName := runtime.FuncForPC(pc).Name()
	return stack{
		funcName: funcName,
		file:     fmt.Sprintf("%s %d", file, line),
	}
}

func usecaseCaller() string {
	pc, _, _, _ := runtime.Caller(2)
	funcName := runtime.FuncForPC(pc).Name()
	return funcName
}
