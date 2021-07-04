package telemux_test

import (
	"fmt"
	"path"
	"reflect"
	"runtime"
	"testing"
)

func assert(condition bool, t *testing.T, arg ...interface{}) {
	if !condition {
		_, file, line, _ := runtime.Caller(1)
		arg = append([]interface{}{fmt.Sprintf("%s:%d:", path.Base(file), line)}, arg...)
		t.Error(arg...)
	}
}

func getFunctionName(i interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
}
