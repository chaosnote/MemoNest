package utils

import "reflect"

func StructToMap(obj any, output map[string]any) {
	v := reflect.ValueOf(obj)
	typeOfS := v.Type()

	for i := 0; i < v.NumField(); i++ {
		fieldName := typeOfS.Field(i).Name
		fieldValue := v.Field(i).Interface()
		output[fieldName] = fieldValue
	}
}
