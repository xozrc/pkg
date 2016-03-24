package reflectutils_test

import (
	"fmt"
	"reflect"
	"testing"
)

import (
	"github.com/stretchr/testify/assert"
	"github.com/xozrc/pkg/reflectutils"
)

var (
	correctDataSlice []interface{}
	errDataSlice     []interface{}
)

func setupTestingData() {
	setupErrTestingData()
	setupCorrectTestingData()
}

func setupErrTestingData() {

	//pointer
	intVal := 1
	intPtr := &intVal
	errDataSlice = append(errDataSlice, intPtr)

	//arr
	arr := []int{1, 1}
	errDataSlice = append(errDataSlice, arr)

	//slice
	sli := make([]int, 0)
	sli = append(sli, 1)
	errDataSlice = append(errDataSlice, sli)

	//map
	m := make(map[string]string, 0)
	m["hello"] = "world"
	errDataSlice = append(errDataSlice, m)
	//chan
	ch := make(chan int, 1)
	errDataSlice = append(errDataSlice, ch)
	//func
	f := func() {}
	errDataSlice = append(errDataSlice, f)

	//todo: example for kind below
	//interace
	//Ptr
	//Struct
	//UnsafePointer

}

func setupCorrectTestingData() {
	correctDataSlice = append(correctDataSlice, false)
	correctDataSlice = append(correctDataSlice, true)
	correctDataSlice = append(correctDataSlice, int(1))
	correctDataSlice = append(correctDataSlice, int8(2))
	correctDataSlice = append(correctDataSlice, int16(3))
	correctDataSlice = append(correctDataSlice, int32(4))
	correctDataSlice = append(correctDataSlice, int64(5))
	correctDataSlice = append(correctDataSlice, uint(1))
	correctDataSlice = append(correctDataSlice, uint8(2))
	correctDataSlice = append(correctDataSlice, uint16(3))
	correctDataSlice = append(correctDataSlice, uint32(4))
	correctDataSlice = append(correctDataSlice, uint64(5))
	correctDataSlice = append(correctDataSlice, float32(1.0))
	correctDataSlice = append(correctDataSlice, float64(2.0))
	correctDataSlice = append(correctDataSlice, "hello")
}

func TestSetPrimitive(t *testing.T) {
	setupTestingData()
	//correct data
	for _, val := range correctDataSlice {
		typ := reflect.TypeOf(val)
		acVal, err := reflectutils.ParsePrimitive(typ, fmt.Sprintf("%v", val))
		assert.NoError(t, err, err)
		assert.EqualValues(t, acVal, val, "correct data failed")
	}

	//no primitive type
	for _, errData := range errDataSlice {
		val := errData
		typ := reflect.TypeOf(val)
		_, err := reflectutils.ParsePrimitive(typ, "")
		_, ok := err.(*reflectutils.InvalidPrimitiveError)
		assert.True(t, ok, "err data failed")
	}

	//parse error
	_, err := reflectutils.ParsePrimitive(reflect.TypeOf(1), "s")
	assert.Error(t, err, "parse err data failed")

	//empty string
	_, err = reflectutils.ParsePrimitive(reflect.TypeOf(1), "")
	assert.Error(t, err, "parse err data failed")
}
