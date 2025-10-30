package xlskit

import (
	"bytes"
	"fmt"

	"github.com/xuri/excelize/v2"
)

type xlskitctrl struct {
	sheetname string
	file      *excelize.File
}

func NewXlsCtrl(sheetname string) *xlskitctrl {
	return &xlskitctrl{
		sheetname: sheetname,
		file:      excelize.NewFile(),
	}
}

func Render(title string, meta []Meta, dataInput []map[string]any) (buf *bytes.Buffer, err error) {
	ctrl := NewXlsCtrl(title)
	return ctrl.Render(meta, dataInput)
}

// meta: 是元数据,{"字段1","字段2","字符3"}
func BuildData(meta []Meta, dataInput []map[string]any) (list [][]CellValue) {

	titleMap := make(map[string]string, 0)
	// 形成字段名称和A,B,C,D之间的映射关系
	fieldArr := []string{}
	for _, value := range meta {
		titleMap[value.Field] = value.Title
		fieldArr = append(fieldArr, value.Field)
	}
	// 添加表头
	rowIndex := 1
	list = make([][]CellValue, 0)
	titlerow := make([]CellValue, 0)
	for index, field := range fieldArr {
		cellName := fmt.Sprintf("%s%d", charfileds[index], rowIndex)
		title := titleMap[field]
		titlerow = append(titlerow, CellValue{
			Name:  cellName,
			Value: title,
		})
	}
	list = append(list, titlerow)

	for rowIndex, dataItem := range dataInput {
		datarow := make([]CellValue, 0)
		for i, field := range fieldArr {
			cellName := fmt.Sprintf("%s%d", charfileds[i], rowIndex+2)
			value := dataItem[field]
			datarow = append(datarow, CellValue{
				Name:  cellName,
				Value: value,
			})
		}
		list = append(list, datarow)
	}
	return list
}

// meta: 是元数据,{"字段1","字段2","字符3"}
func (s *xlskitctrl) Render(meta []Meta, dataInput []map[string]any) (buf *bytes.Buffer, err error) {
	index, err := s.file.NewSheet(s.sheetname)
	if err != nil {
		return
	}
	s.file.SetActiveSheet(index)
	data := BuildData(meta, dataInput)
	for _, row := range data {
		for _, c := range row {
			s.file.SetCellValue(s.sheetname, c.Name, c.Value)
		}
	}
	return s.file.WriteToBuffer()
}
