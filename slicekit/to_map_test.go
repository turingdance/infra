package slicekit

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestToMap(t *testing.T) {
	type args struct {
		slice     []int
		keyFunc   func(int, int, []int) int
		valueFunc func(int, int, []int) any
	}
	tests := []struct {
		name string
		args args
		want map[int]any
	}{
		{
			name: "number1",
			args: args{
				slice: []int{2, 3, 4, 13},
				keyFunc: func(item int, index int, slice []int) int {
					return item
				},
				valueFunc: func(item int, index int, slice []int) any {
					return item
				},
			},
			want: map[int]any{
				2: 2, 3: 3, 4: 4, 13: 13,
			},
		},
		{
			name: "number2",
			args: args{
				slice: []int{2, 3, 4, 13},
				keyFunc: func(item int, index int, slice []int) int {
					return item
				},
				valueFunc: func(item int, index int, slice []int) any {
					return fmt.Sprint(item)
				},
			},
			want: map[int]any{
				2: "2", 3: "3", 4: "4", 13: "13",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ToMap(tt.args.slice, tt.args.keyFunc, tt.args.valueFunc)
			gotJson, _ := json.Marshal(&got)
			wantJson, _ := json.Marshal(&tt.want)
			if string(gotJson) != string(wantJson) {
				t.Errorf("ToMap() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestToMapList(t *testing.T) {
	type User struct {
		Name  string `json:"name"`
		Age   int    `json:"age,omitempty"`
		Email string `json:"email_address"`
	}

	users := []User{
		{Name: "Alice", Age: 30, Email: "alice@example.com"},
		{Name: "Bob", Age: 0, Email: "bob@example.com"},
	}

	maps := ObjListToMapList(users)
	for i, v := range maps {
		t.Run(v["name"].(string), func(t *testing.T) {
			if v["name"] != users[i].Name {
				t.Error(v)
			}
		})
	}
}
