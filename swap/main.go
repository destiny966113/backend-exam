package main

import (
	"fmt"
	"reflect"
)

func swap[T any](a, b T) {
	// 檢查傳入的參數是否為指標
	aVal := reflect.ValueOf(a)
	bVal := reflect.ValueOf(b)

	if aVal.Kind() != reflect.Ptr || bVal.Kind() != reflect.Ptr {
		panic("swap: arguments must be pointers")
	}

	// 檢查指標是否為 nil
	if aVal.IsNil() || bVal.IsNil() {
		panic("swap: nil pointer")
	}

	// 取得指標指向的值
	aElem := aVal.Elem()
	bElem := bVal.Elem()

	// 檢查是否可以設定值
	if !aElem.CanSet() || !bElem.CanSet() {
		panic("swap: cannot set value")
	}

	// 建立臨時變數儲存 a 的值
	temp := reflect.New(aElem.Type()).Elem()
	temp.Set(aElem)

	// 交換值
	aElem.Set(bElem)
	bElem.Set(temp)
}

func main() {
	a := 10
	b := 20

	fmt.Printf("a = %d, &a = %p\n", a, &a)
	fmt.Printf("b = %d, &b = %p\n", b, &b)

	swap(&a, &b)

	fmt.Printf("a = %d, &a = %p\n", a, &a)
	fmt.Printf("b = %d, &b = %p\n", b, &b)
}
