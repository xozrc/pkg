package reflectutils

import (
	"fmt"
	"reflect"
	"strconv"
)

type InvalidPrimitiveError struct {
	typ reflect.Type
}

func (ipe *InvalidPrimitiveError) Error() string {
	return fmt.Sprintf("invalid primitive type (%s)", ipe.typ.String())
}

//todo:improve complex
//except complex64 and complex 128
func ParsePrimitive(typ reflect.Type, val string) (interface{}, error) {

	tvalue := reflect.New(typ)
	value := tvalue.Elem()
	switch typ.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:

		if val == "" {
			val = "0"
		}
		intVal, err := strconv.ParseInt(val, 10, 64)
		if err != nil {
			return nil, err
		}
		value.SetInt(intVal)
		return value.Interface(), nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:

		if val == "" {
			val = "0"
		}
		uintVal, err := strconv.ParseUint(val, 10, 64)
		if err != nil {
			return nil, err
		}
		value.SetUint(uintVal)
		return value.Interface(), nil
	case reflect.Bool:

		if val == "" {
			val = "false"
		}
		boolVal, err := strconv.ParseBool(val)
		if err != nil {
			return nil, err
		}

		value.SetBool(boolVal)
		return value.Interface(), nil
	case reflect.Float32:

		if val == "" {
			val = "0.0"
		}
		floatVal, err := strconv.ParseFloat(val, 32)
		if err != nil {
			return nil, err
		}
		value.SetFloat(floatVal)
		return value.Interface(), nil
	case reflect.Float64:

		if val == "" {
			val = "0.0"
		}
		floatVal, err := strconv.ParseFloat(val, 64)
		if err != nil {
			return nil, err
		}
		value.SetFloat(floatVal)
		return value.Interface(), nil
	case reflect.String:

		value.SetString(val)
		return value.Interface(), nil
	default:
		return nil, &InvalidPrimitiveError{typ}
	}
}
