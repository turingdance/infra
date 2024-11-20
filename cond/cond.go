package cond

import (
	"errors"
	"fmt"
	"strings"
	"unicode"
)

type Cond struct {
	Field    string      `json:"field"`
	DataType Datatype    `json:"datatype"`
	Op       OPTYPE      `json:"op"`
	Value    interface{} `json:"value"`
	KeyFunc  KeyFuncType `json:"keyfunc"`
}

func (c Cond) GetField() string {
	if c.KeyFunc == LowerCamel {
		return toLowerCamelCase(c.Field)
	} else if c.KeyFunc == UpperCamel {
		return toUpperCamelCase(c.Field)
	} else if c.KeyFunc == SnakeCase {
		return toSnakeCase(c.Field)
	} else {
		return toSnakeCase(c.Field)
	}
}

const PAGESIZEMAX = 1024 * 1024

// 转换
func (c Cond) Build() (sql string, value interface{}, err error) {
	err = nil
	if c.Op == OPEQ {
		return fmt.Sprintf("%s = ?", c.GetField()), c.Value, err
	} else if c.Op == OPEGT {
		return fmt.Sprintf("%s >= ?", c.GetField()), c.Value, err
	} else if c.Op == OPGT {
		return fmt.Sprintf("%s > ?", c.GetField()), c.Value, err
	} else if c.Op == OPLT {
		return fmt.Sprintf("%s < ?", c.GetField()), c.Value, err
	} else if c.Op == OPLET {
		return fmt.Sprintf("%s <= ?", c.GetField()), c.Value, err
	} else if c.Op == OPLIKE {
		return fmt.Sprintf("%s like ?", c.GetField()), "%" + fmt.Sprintf("%s", c.Value) + "%", err
	} else if c.Op == OPIN {
		return fmt.Sprintf("%s in ?", c.GetField()), c.Value, err
	} else if c.Op == OPBETWEEN {
		if c.DataType == INTARR {
			value = 1
			err = nil
			switch c.Value.(type) {
			case []int:
				sql = fmt.Sprintf("%s between %d and %d and 1 = ?", c.GetField(), c.Value.([]int)[0], c.Value.([]int)[1])
			case []int16:
				sql = fmt.Sprintf("%s between %d and %d and 1 = ?", c.GetField(), c.Value.([]int16)[0], c.Value.([]int16)[1])
			case []int32:
				sql = fmt.Sprintf("%s between %d and %d and 1 = ?", c.GetField(), c.Value.([]int32)[0], c.Value.([]int32)[1])
			case []int64:
				sql = fmt.Sprintf("%s between %d and %d and 1 = ?", c.GetField(), c.Value.([]int64)[0], c.Value.([]int64)[1])
			case []uint:
				sql = fmt.Sprintf("%s between %d and %d and 1 = ?", c.GetField(), c.Value.([]uint)[0], c.Value.([]uint)[1])
			case []float32:
				sql = fmt.Sprintf("%s between %0f and %0f and 1 = ?", c.GetField(), c.Value.([]float32)[0], c.Value.([]float32)[1])
			case []float64:
				sql = fmt.Sprintf("%s between %0f and %0f and 1 = ?", c.GetField(), c.Value.([]float64)[0], c.Value.([]float64)[1])
			default:
				sql = ""
				err = errors.New("数据格式不正确")
			}
			return sql, value, err
		} else if c.DataType == STRARR {
			if value, is := c.Value.([]string); is {
				return fmt.Sprintf("%s between %s and %s and 1 = ?", c.GetField(), value[0], value[1]), 1, err
			} else {
				return "", "", errors.New("数据格式不正确")
			}

		} else {
			if value, is := c.Value.([]string); is {
				return fmt.Sprintf("%s between %s and %s and 1 = ?", c.GetField(), value[0], value[1]), 1, err
			} else if value, is := c.Value.([]int); is {
				return fmt.Sprintf("%s between %d and %d and 1 = ?", c.GetField(), value[0], value[1]), 1, err
			} else {
				return " 1 = ? ", 1, err
			}
		}
	} else {
		return " 1 = ? ", 1, err
	}
}

type KeyFunc func(str string) string
type KeyFuncType string

// toLowerCamelCase 将字符串转换为小驼峰格式
func toLowerCamelCase(s string) string {
	words := splitByInitialisms(s)
	var camelCaseStr strings.Builder
	camelCaseStr.Grow(len(s))

	for i, word := range words {
		if i == 0 {
			camelCaseStr.WriteString(strings.ToLower(word))
		} else {
			camelCaseStr.WriteString(strings.Title(strings.ToLower(word)))
		}
	}

	return camelCaseStr.String()
}

// toUpperCamelCase 将字符串转换为大驼峰格式
func toUpperCamelCase(s string) string {
	words := splitByInitialisms(s)
	var camelCaseStr strings.Builder
	camelCaseStr.Grow(len(s))

	for i, word := range words {
		if i == 0 {
			camelCaseStr.WriteString(strings.ToUpper(word))
		} else {
			camelCaseStr.WriteString(strings.Title(strings.ToLower(word)))
		}
	}

	return camelCaseStr.String()
}

// toSnakeCase 将字符串转换为蛇形（下划线）格式
func toSnakeCase(s string) string {
	words := splitByInitialisms(s)
	var snakeCaseStr strings.Builder
	snakeCaseStr.Grow(len(s) + 2*len(words) - 1) // 加上下划线

	for i, word := range words {
		if i > 0 {
			snakeCaseStr.WriteString("_")
		}
		snakeCaseStr.WriteString(strings.ToLower(word))
	}

	return snakeCaseStr.String()
}

// splitByInitialisms 分割字符串，考虑首字母大写和缩写词
func splitByInitialisms(s string) []string {
	var words []string
	var word strings.Builder

	for _, ch := range s {
		if ch == '_' {
			if word.Len() > 0 {
				words = append(words, word.String())
				word.Reset()
			}
		} else if ch == ' ' {
			if word.Len() > 0 {
				words = append(words, word.String())
				word.Reset()
			}
		} else if unicode.IsUpper(ch) {
			if word.Len() > 0 {
				words = append(words, word.String())
				word.Reset()
			}
			word.WriteRune(unicode.ToLower(ch))
		} else {
			word.WriteRune(ch)
		}
	}

	if word.Len() > 0 {
		words = append(words, word.String())
	}

	return words
}

const (
	LowerCamel KeyFuncType = "lowercamel"
	UpperCamel KeyFuncType = "uppercamel"
	SnakeCase  KeyFuncType = "snake"
)

type CondWraper struct {
	Pager   Pager       `json:"pager"`
	Order   Order       `json:"order"`
	Conds   []Cond      `json:"conds"`
	KeyFunc KeyFuncType `json:"keyfunc"`
}

func NewListAllWraper() *CondWraper {
	return &CondWraper{
		Conds:   make([]Cond, 0),
		Order:   Order{},
		Pager:   Pager{Pagefrom: 0, Pagesize: PAGESIZEMAX},
		KeyFunc: SnakeCase,
	}
}
func NewCondWrapper() *CondWraper {
	return &CondWraper{
		Conds:   make([]Cond, 0),
		Order:   Order{},
		Pager:   Pager{Pagefrom: 0, Pagesize: 20},
		KeyFunc: SnakeCase,
	}
}
func (c *CondWraper) AddCond(conds ...Cond) *CondWraper {
	for _, v := range conds {
		if v.KeyFunc == "" {
			v.KeyFunc = SnakeCase
		}
		c.Conds = append(c.Conds, v)
	}
	return c
}
func (c *CondWraper) AddOneCond(field string, op OPTYPE, value any) *CondWraper {
	c.Conds = append(c.Conds, Cond{
		Field: field, Op: op, Value: value,
	})
	return c
}
func (c *CondWraper) SetPager(pager Pager) *CondWraper {
	c.Pager = pager
	return c
}

func (c *CondWraper) SetKeyFunc(keyfun KeyFuncType) *CondWraper {
	c.KeyFunc = keyfun
	return c
}

func (c *CondWraper) Pagesize(size int) *CondWraper {
	c.Pager.Pagesize = size
	return c
}
func (c *CondWraper) Pagefrom(from int) *CondWraper {
	c.Pager.Pagefrom = from
	return c
}

func (c *CondWraper) SetOrer(order Order) *CondWraper {
	c.Order = order
	return c
}
