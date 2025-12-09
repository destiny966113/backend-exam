package main

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
)

func TrimAllStrings(a any) {
	visited := make(map[uintptr]bool)
	trimAllStringsRecursive(reflect.ValueOf(a), visited)
}

func trimAllStringsRecursive(v reflect.Value, visited map[uintptr]bool) {
	// 處理無效值
	if !v.IsValid() {
		return
	}

	// 處理指標類型
	if v.Kind() == reflect.Ptr {
		if v.IsNil() {
			return
		}

		// 檢查是否已訪問過此指標（避免循環引用）
		ptr := v.Pointer()
		if visited[ptr] {
			return
		}
		visited[ptr] = true

		// 遞迴處理指標指向的值
		trimAllStringsRecursive(v.Elem(), visited)
		return
	}

	// 處理介面類型
	if v.Kind() == reflect.Interface {
		if v.IsNil() {
			return
		}
		trimAllStringsRecursive(v.Elem(), visited)
		return
	}

	// 處理結構體類型
	if v.Kind() == reflect.Struct {
		for i := 0; i < v.NumField(); i++ {
			field := v.Field(i)
			if field.CanSet() {
				trimAllStringsRecursive(field, visited)
			}
		}
		return
	}

	// 處理字串類型
	if v.Kind() == reflect.String {
		if v.CanSet() {
			trimmed := strings.TrimSpace(v.String())
			v.SetString(trimmed)
		}
		return
	}

	// 處理切片類型
	if v.Kind() == reflect.Slice {
		for i := 0; i < v.Len(); i++ {
			trimAllStringsRecursive(v.Index(i), visited)
		}
		return
	}

	// 處理陣列類型
	if v.Kind() == reflect.Array {
		for i := 0; i < v.Len(); i++ {
			trimAllStringsRecursive(v.Index(i), visited)
		}
		return
	}

	// 處理 map 類型
	if v.Kind() == reflect.Map {
		for _, key := range v.MapKeys() {
			val := v.MapIndex(key)
			// 對於 map，我們需要建立新的值來設定
			newVal := reflect.New(val.Type()).Elem()
			newVal.Set(val)
			trimAllStringsRecursive(newVal, visited)
			v.SetMapIndex(key, newVal)
		}
		return
	}
}

func main() {
	type Person struct {
		Name string
		Age  int
		Next *Person
	}

	a := &Person{
		Name: " name ",
		Age:  20,
		Next: &Person{
			Name: " name2 ",
			Age:  21,
			Next: &Person{
				Name: " name3 ",
				Age:  22,
			},
		},
	}

	TrimAllStrings(&a)

	m, _ := json.Marshal(a)

	fmt.Println(string(m))

	a.Next = a

	TrimAllStrings(&a)

	fmt.Println(a.Next.Next.Name == "name")
}
