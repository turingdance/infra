package treekit

import (
	"fmt"
	"reflect"
)

// TreeNode 表示树的节点结构
type TreeNode struct {
	Data     interface{} `json:"data"`
	Children []*TreeNode `json:"children"`
}

// ToTree 将任意结构图数组转换为树形结构
// slice 是包含结构体的切片，每个结构体必须有 ID 和 ParentID 字段
// idField 和 parentIDField 分别是 ID 和 ParentID 字段的名称
func ToTree(slice interface{}, idField, parentIDField string) ([]*TreeNode, error) {
	// 获取反射值和类型
	value := reflect.ValueOf(slice)
	if value.Kind() != reflect.Slice {
		return nil, fmt.Errorf("输入不是一个切片")
	}

	// 检查元素类型是否为结构体
	if value.Type().Elem().Kind() != reflect.Struct {
		return nil, fmt.Errorf("切片元素不是结构体")
	}

	// 创建节点映射
	nodes := make(map[interface{}]*TreeNode)
	rootNodes := make([]*TreeNode, 0)

	// 构建节点映射
	for i := 0; i < value.Len(); i++ {
		item := value.Index(i)
		id := getFieldValue(item, idField)
		parentID := getFieldValue(item, parentIDField)

		if id == nil {
			return nil, fmt.Errorf("第 %d 个元素的 %s 字段不存在", i, idField)
		}

		if parentID == nil {
			return nil, fmt.Errorf("第 %d 个元素的 %s 字段不存在", i, parentIDField)
		}

		// 创建节点
		node := &TreeNode{
			Data:     item.Interface(),
			Children: []*TreeNode{},
		}
		nodes[id] = node

		// 如果是根节点
		if parentID == reflect.Zero(reflect.TypeOf(parentID)).Interface() {
			rootNodes = append(rootNodes, node)
		}
	}

	// 构建树结构
	for i := 0; i < value.Len(); i++ {
		item := value.Index(i)
		id := getFieldValue(item, idField)
		parentID := getFieldValue(item, parentIDField)

		// 如果不是根节点，添加到父节点的子节点列表中
		if parentID != reflect.Zero(reflect.TypeOf(parentID)).Interface() {
			parentNode, exists := nodes[parentID]
			if exists {
				if childNode, ok := nodes[id]; ok {
					parentNode.Children = append(parentNode.Children, childNode)
				}
			}
		}
	}

	return rootNodes, nil
}

// getFieldValue 通过反射获取结构体字段的值
func getFieldValue(structValue reflect.Value, fieldName string) interface{} {
	// 如果是指针，解引用
	if structValue.Kind() == reflect.Ptr {
		structValue = structValue.Elem()
	}

	// 获取字段
	field := structValue.FieldByName(fieldName)
	if !field.IsValid() {
		return nil
	}

	// 获取字段值
	return field.Interface()
}
