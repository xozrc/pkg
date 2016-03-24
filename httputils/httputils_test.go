package httputils_test

import (
	"fmt"
	"reflect"
	"strings"
	"testing"
)

import (
	"github.com/stretchr/testify/assert"
	"github.com/xozrc/pkg/httputils"
)

var (
	correctTestFormSlice []*testForm
	errotTestFormSlice   []*errorTestForm
)

type errorTestForm struct {
	TestString string   `form:"testS"`
	TestBool   bool     `form:"testB"`
	TestIntPtr *int     `form:"testI"`
	TestUint   uint     `form:"testU"`
	TestSlice  []string `form:"testSlice"`
}

type testForm struct {
	TestString string   `form:"testS"`
	TestBool   bool     `form:"testB"`
	TestInt    int      `form:"testI"`
	TestUint   uint     `form:"testU"`
	TestSlice  []string `form:"testSlice"`
}

func (tf *testForm) Equal(tf2 *testForm) bool {
	if tf2 == nil {
		return false
	}
	if !strings.EqualFold(tf.TestString, tf2.TestString) {
		return false
	}
	if tf.TestBool != tf2.TestBool {
		return false
	}
	if tf.TestInt != tf2.TestInt {
		return false
	}
	if tf.TestUint != tf2.TestUint {
		return false
	}

	tlen := len(tf.TestSlice)
	if len(tf2.TestSlice) != tlen {
		return false
	}
	for i := 0; i < tlen; i++ {
		if !strings.EqualFold(tf.TestSlice[i], tf2.TestSlice[i]) {
			return false
		}
	}
	return true
}

func convertToMap(obj interface{}) map[string][]string {
	result := make(map[string][]string, 0)
	typ := reflect.TypeOf(obj)
	val := reflect.ValueOf(obj)
	typ = typ.Elem()
	val = val.Elem()

	numFields := val.NumField()
	for i := 0; i < numFields; i++ {
		fieldTyp := typ.Field(i)
		fieldVal := val.Field(i)

		key := fieldTyp.Tag.Get("form")
		tmpSlice := make([]string, 0)
		if fieldTyp.Type.Kind() == reflect.Slice {
			for j := 0; j < fieldVal.Len(); j++ {
				tmpSlice = append(tmpSlice, fmt.Sprintf("%v", fieldVal.Index(j).Interface()))
			}
		} else {
			tmpSlice = append(tmpSlice, fmt.Sprintf("%v", fieldVal.Interface()))

		}
		result[key] = tmpSlice
	}
	return result
}

func TestBindForm(t *testing.T) {
	testCorrectBindForm(t)
}

func testCorrectBindForm(t *testing.T) {
	tff := &testForm{}
	tff.TestInt = 1
	tff.TestBool = false
	tff.TestString = "test"
	tff.TestUint = uint(1)
	tff.TestSlice = []string{"test1", "test2"}
	tf := &testForm{}
	err := httputils.BindForm(tf, convertToMap(tff), nil)
	assert.NoError(t, err, "bind error")
	assert.True(t, tf.Equal(tff), "error")
}
