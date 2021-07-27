package telemux_test

import (
	"io/ioutil"
	"os"
	"reflect"
	"testing"

	tm "github.com/and3rson/telemux/v2"
)

func TestPersistenceKey(t *testing.T) {
	p := tm.PersistenceKey{"foo", 13, 37}
	text, _ := p.MarshalText()
	assert(string(text) == "foo:13:37", t)

	assert(p.UnmarshalText([]byte("bar:42:69")) == nil, t)
	assert(reflect.DeepEqual(p, tm.PersistenceKey{"bar", 42, 69}), t)

	assert(p.UnmarshalText([]byte("bar:42:bb")) != nil, t)
	assert(p.UnmarshalText([]byte("bar:aa:69")) != nil, t)
}

func TestFilePersistence(t *testing.T) {
	f, err := ioutil.TempFile("", "telemux_persistence")
	if err != nil {
		t.Error("Failed to create temporary file")
	}
	defer func() {
		f.Close()
		os.Remove(f.Name())
	}()

	f.Close()
	os.Remove(f.Name())

	pk1 := tm.PersistenceKey{"foo", 1, 2}
	pk2 := tm.PersistenceKey{"bar", 3, 4}

	p := tm.NewFilePersistence(f.Name())
	assert(reflect.DeepEqual(p.GetData(pk1), tm.Data{}), t)
	assert(reflect.DeepEqual(p.GetData(pk2), tm.Data{}), t)

	p.SetData(pk1, tm.Data{"foo": "bar", "number": 69})
	assert(reflect.DeepEqual(p.GetData(pk1), tm.Data{"foo": "bar", "number": 69.0}), t)
	assert(reflect.DeepEqual(p.GetData(pk2), tm.Data{}), t)

	assert(p.GetState(pk1) == "", t)
	p.SetState(pk2, "state2")
	assert(p.GetState(pk1) == "", t)
	assert(p.GetState(pk2) == "state2", t)
}
