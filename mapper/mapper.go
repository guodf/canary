package mapper

import (
	"reflect"
)

// Mapper mapper
func Mapper(from, to interface{}) {
	analysisType(from)
	analysisType(to)
	toElem := reflect.ValueOf(to)
	if toElem.Type().Kind() != reflect.Ptr {
		panic("to必须是指针类型")
	}

	if toElem.IsZero() {
		panic("to没有明确的地址")
	}

	fromElem := reflect.ValueOf(from)
	if fromElem.Kind() == reflect.Ptr {
		fromElem = fromElem.Elem()
	}
	if fields, ok := getFields(to); ok {
		for _, field := range fields {
			toField := toElem.Elem().FieldByName(field.Name)
			setValue(toField, fromElem.FieldByName(field.Name))
		}
	}
}

func setValue(toField, fromField reflect.Value) {
	switch toField.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		toField.SetInt(fromField.Int())
		break
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		toField.SetUint(fromField.Uint())
		break
	case reflect.String:
		toField.SetString(fromField.String())
		break
		// case reflect.Struct:
		// 	toField.Set(fromField)
		// 	break
	}
}

func getFields(obj interface{}) ([]reflect.StructField, bool) {
	key := reflect.TypeOf(obj).Elem().Name()

	elems, ok := structs[key]
	return elems, ok
}

var structs map[string][]reflect.StructField = make(map[string][]reflect.StructField, 0)

func analysisType(arg interface{}) {
	if arg == nil {
		return
	}

	typ, ok := arg.(reflect.Type)
	if !ok {
		typ = reflect.TypeOf(arg)
	}
	switch typ.Kind() {
	case reflect.Ptr:
		analysisStruct(typ.Elem())
		break
	case reflect.Struct:
		analysisStruct(typ)
		break
	case reflect.Int, reflect.String:
		break
	default:
		break
	}
}

func analysisStruct(typ reflect.Type) {
	fields := []reflect.StructField{}
	structName := typ.Name()
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		if len(field.PkgPath) > 0 {
			continue
		}
		fields = append(fields, field)
		analysisType(field.Type)
	}

	if elems, ok := structs[structName]; ok {
		elems = append(elems, fields...)
		return
	}
	structs[structName] = fields
}
