package telemux_test

import "testing"

func assert(condition bool, t *testing.T, arg ...interface{}) {
	if !condition {
		t.Error(arg...)
	}
}
