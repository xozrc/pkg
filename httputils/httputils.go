package httputils

import (
	"errors"
	"net/http"
	"reflect"
	"strings"
)

const (
	formKey string = "form"
)

func Bind(req *http.Request, obj interface{}) error {
	if req == nil {
		return errors.New("request is nil")
	}

	contentType := req.Header.Get("Content-Type")

	if strings.Contains(contentType, "form-urlencoded") {
		return BindForm(req, obj)
	} else if strings.Contains(contentType, "multipart/form-data") {

	} else if strings.Contains(contentType, "application/json") {
		return JsonForm(req, obj)
	} else {
		if contentType == "" {
			return errors.New("empty Content-Type")
		} else {
			return errors.New("Unsupported Content-Type")
		}

	}
	return nil
}

func BindForm(req *http.Request, form interface{}) error {
	err := req.ParseForm()
	if err != nil {
		return err
	}
	mapForm(form, req.Form, nil)

}

func mapForm(form interface{}, form map[string][]string) (err error) {

	formVal := reflect.ValueOf(form)
	typ := reflect.TypeOf(form)

	if formVal.Kind() != reflect.Ptr || formVal.IsNil() {
		err = NewInvalidBindError(typ)
		return
	}

	for i := 0; i < typ.NumField(); i++ {
		fieldType := typ.Field(i)
		fieldVal := formVal.Field(i)
		if tagName := fieldType.Tag.Get(formKey); tagName != "" {
			if !fieldVal.CanSet() {
				continue
			}

			var tagVal []string
			if tagVal, ok := form[tagName]; !ok {
				if tagVal, ok = formfile[tagName]; !ok {
					continue
				}
			}

			if len(tagVal) == 0 {
				fieldVal.Set(reflect.Zero(fieldType))
				continue
			}

			switch fieldType.Type.Kind() {
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:

				val := tagVal[0]
				if val == "" {
					val = "0"
				}
				intVal, err := strconv.ParseInt(val, 10, 64)
				if err != nil {
					return err
				}
				fieldVal.SetInt(intVal)

			case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
				val := tagVal[0]
				if val == "" {
					val = "0"
				}
				uintVal, err := strconv.ParseUint(val, 10, 64)
				if err != nil {
					return err
				}
				fieldVal.SetUint(uintVal)

			case reflect.Bool:
				val := tagVal[0]
				if val == "" {
					val = "false"
				}
				boolVal, err := strconv.ParseBool(val)
				if err != nil {
					return err
				}
				fieldVal.SetBool(boolVal)

			case reflect.Float32:
				val := tagVal[0]
				if val == "" {
					val = "0.0"
				}
				floatVal, err := strconv.ParseFloat(val, 32)
				if err != nil {
					return err
				}
				fieldVal.SetFloat(floatVal)

			case reflect.Float64:
				val := tagVal[0]
				if val == "" {
					val = "0.0"
				}
				floatVal, err := strconv.ParseFloat(val, 64)
				if err != nil {
					return err
				}
				fieldVal.SetFloat(floatVal)

			case reflect.String:
				val := tagVal[0]
				fieldVal.SetString(val)
			}

		}

		switch fieldType.Type.Kind() {
		case reflect.Interface:
		case reflect.Ptr:
		case reflect.Struct:
		case reflect.Map:
		case reflect.Array:
		case reflect.Slice:
			return NewInvalidBindFieldError(fieldType.Type)
			break
		default:
			if tagName := fieldType.Tag.Get(formKey); tagName != "" {
				if !fieldVal.CanSet() {
					continue
				}
				var tagVal
				if tagVal, ok := form[tagName]; !ok {
					if tagVal, ok = formfile[tagName]; !ok {
						continue
					}
				}

			}
			break
		}
		if err != nil {
			return
		}

	}

	// for i := 0; i < typ.NumField(); i++ {
	// 	typeField := typ.Field(i)
	// 	structField := formStruct.Field(i)

	// 	if typeField.Type.Kind() == reflect.Ptr && typeField.Anonymous {
	// 		structField.Set(reflect.New(typeField.Type.Elem()))
	// 		mapForm(structField.Elem(), form, formfile, errors)
	// 		if reflect.DeepEqual(structField.Elem().Interface(), reflect.Zero(structField.Elem().Type()).Interface()) {
	// 			structField.Set(reflect.Zero(structField.Type()))
	// 		}
	// 	} else if typeField.Type.Kind() == reflect.Struct {
	// 		mapForm(structField, form, formfile, errors)
	// 	} else if inputFieldName := typeField.Tag.Get("form"); inputFieldName != "" {
	// 		if !structField.CanSet() {
	// 			continue
	// 		}

	// 		inputValue, exists := form[inputFieldName]
	// 		if exists {
	// 			numElems := len(inputValue)
	// 			if structField.Kind() == reflect.Slice && numElems > 0 {
	// 				sliceOf := structField.Type().Elem().Kind()
	// 				slice := reflect.MakeSlice(structField.Type(), numElems, numElems)
	// 				for i := 0; i < numElems; i++ {
	// 					setWithProperType(sliceOf, inputValue[i], slice.Index(i), inputFieldName, errors)
	// 				}
	// 				formStruct.Field(i).Set(slice)
	// 			} else {
	// 				setWithProperType(typeField.Type.Kind(), inputValue[0], structField, inputFieldName, errors)
	// 			}
	// 			continue
	// 		}

	// 		inputFile, exists := formfile[inputFieldName]
	// 		if !exists {
	// 			continue
	// 		}
	// 		fhType := reflect.TypeOf((*multipart.FileHeader)(nil))
	// 		numElems := len(inputFile)
	// 		if structField.Kind() == reflect.Slice && numElems > 0 && structField.Type().Elem() == fhType {
	// 			slice := reflect.MakeSlice(structField.Type(), numElems, numElems)
	// 			for i := 0; i < numElems; i++ {
	// 				slice.Index(i).Set(reflect.ValueOf(inputFile[i]))
	// 			}
	// 			structField.Set(slice)
	// 		} else if structField.Type() == fhType {
	// 			structField.Set(reflect.ValueOf(inputFile[0]))
	// 		}
	// 	}
	// }
}

func WriteJSON(w http.ResponseWriter, code int, v interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	return json.NewEncoder(w).Encode(v)
}
