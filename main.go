package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// 因为要导出，所以这里属性名一定要大写
// `json:"name"`之间不能带空格，不然会报错且会致使json转小写别名无效
type Person struct {
	Name string `json:"name"`
	Age  uint8  `json:"age"`
	Id   int    `json:"id"`
}

type Res struct {
	Code string `json:"code"`
	Data Person `json:"data"`
}

func Index(w http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	person := Person{
		Name: "dilireba",
		Age:  18,
		Id:   202301,
	}
	res := Res{
		Code: "0",
		Data: person,
	}
	jsonData, _ := json.Marshal(res)
	w.Write(jsonData)
}

type AddRes struct {
	Code    string `json:"code"`
	Data    bool   `json:"data"`
	Message string `json:"message"`
}

func Add(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

	user := r.FormValue("user")
	fmt.Printf("user: %s\n", user)
	if user == "" {
		res := AddRes{
			Code:    "1",
			Data:    false,
			Message: "缺少必要参数：user",
		}

		jsonData, _ := json.Marshal(res)
		w.Write(jsonData)
		return
	}

	res := AddRes{
		Code:    "0",
		Data:    true,
		Message: fmt.Sprintf("新增了用户：%s", user),
	}

	jsonData, _ := json.Marshal(res)
	w.Write(jsonData)
}

func main() {
	router := httprouter.New()

	// 预检请求
	router.GlobalOPTIONS = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Access-Control-Request-Method") != "" {
			// Set CORS headers
			header := w.Header()
			header.Set("Access-Control-Allow-Methods", header.Get("Allow"))
			header.Set("Access-Control-Allow-Origin", "*")
		}

		// Adjust status code to 204
		w.WriteHeader(http.StatusNoContent)
	})

	router.GET("/", Index)
	router.POST("/add", Add)

	fmt.Println("server is running at 5501...")
	log.Fatal(http.ListenAndServe(":5501", router))
}
