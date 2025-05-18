package handler

import (
	"io"
	"net/http"
	"os"
)

// 处理文件上传
func UpLoadHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		// 处理GET请求
		data, err := os.ReadFile("./static/view/index.html")
		if err != nil {
			io.WriteString(w, "internel server error")
			return
		}
		io.WriteString(w, string(data))
	case "POST":
		// 处理POST请求
	}
}
