package slicekit

import (
	"reflect"
	"strings"
)

// ToMap iterates the slice to map following the keyFunc and valueFunc
func ToMap[T any, K comparable, V any](slice []T, keyFunc func(item T, index int, slice []T) K,
	valueFunc func(item T, index int, slice []T) V) map[K]V {
	dict := make(map[K]V, len(slice))
	for index, value := range slice {
		dict[keyFunc(value, index, slice)] = valueFunc(value, index, slice)
	}
	return dict
}

// ObjListToMapList 将结构体切片转换为map切片
func ObjListToMapList[T any](rows []T) []map[string]any {
	result := make([]map[string]any, len(rows))

	for i, row := range rows {
		rowMap := make(map[string]any)

		// 获取反射值和类型
		v := reflect.ValueOf(row)

		// 处理指针
		if v.Kind() == reflect.Ptr {
			if v.IsNil() {
				continue
			}
			v = v.Elem()
		}

		// 确保是结构体
		if v.Kind() != reflect.Struct {
			continue
		}

		// 获取类型
		t := v.Type()

		// 遍历结构体字段
		for j := 0; j < v.NumField(); j++ {
			field := t.Field(j)

			// 处理未导出字段
			if !v.Field(j).CanInterface() {
				continue
			}

			// 获取字段名（优先使用json标签，否则使用字段名）
			fieldName := field.Tag.Get("json")
			if fieldName == "" {
				fieldName = field.Name
			}

			// 处理标签中的选项（如omitempty）
			if idx := strings.Index(fieldName, ","); idx != -1 {
				fieldName = fieldName[:idx]
			}

			// 跳过空字段（如果有omitempty标签）
			if strings.Contains(field.Tag.Get("json"), "omitempty") {
				if isEmptyValue(v.Field(j)) {
					continue
				}
			}

			// 获取字段值
			value := v.Field(j).Interface()

			// 添加到map
			rowMap[fieldName] = value
		}

		result[i] = rowMap
	}

	return result
}

// isEmptyValue 判断值是否为空（用于omitempty）
func isEmptyValue(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Array, reflect.Map, reflect.Slice, reflect.String:
		return v.Len() == 0
	case reflect.Bool:
		return !v.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return v.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0
	case reflect.Interface, reflect.Ptr:
		return v.IsNil()
	}
	return false
}
