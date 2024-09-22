package common

import (
	"encoding/json"
	"fmt"
)

// 打印一个结构体的json字符串
func JsonString(obj any) {
	js, err := json.Marshal(obj)
	if err != nil {
		fmt.Println("Can not convernt to json!")
		return
	}
	fmt.Println((string(js)))
}
