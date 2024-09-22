package router_test

import (
	"ginEssential/controller/router"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

// import "ginEssential/ocean_learn/model"

func TestHelloHandler(t *testing.T) {
	// 创建一个新的请求
	req, err := http.NewRequest("GET", "/hello", nil)
	if err != nil {
		t.Fatal(err)
	}

	// 创建一个响应记录器
	rr := httptest.NewRecorder()

	// 对路由使用了 ServeHTTP 方法，所以这里调用它
	r := router.Router()
	r.ServeHTTP(rr, req)

	// 检查返回的状态码
	assert.Equal(t, http.StatusOK, rr.Code)
	// 或者这么写
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code. got %v want %v",
			status, http.StatusOK)
	}

	// 检查返回的内容
	// body, err := ioutil.ReadAll(rr.Result().Body)
	// if err != nil {
	// 	t.Errorf("Error reading body: %v", err)
	// }

	// expected := "<html><body><h1>Hello, world!</h1></body></html>\n"
	// if string(body) != expected {
	// 	t.Errorf("handler returned unexpected body: got %v want %v", body, expected)
	// }
}

// 测试注册
func TestRegister(t *testing.T) {
}
