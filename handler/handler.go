package handler

import (
	"YunDisk/meta"
	"YunDisk/util"
	"io"
	"net/http"
	"os"
	"time"
)

// 处理文件上传
func UpLoadHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		// 处理GET请求
		data, err := os.ReadFile("static/view/index.html")
		if err != nil {
			io.WriteString(w, "internel server error")
			return
		}
		io.WriteString(w, string(data))
	case "POST":
		// 处理POST请求
		r.FormFile("file")
		file, header, err := r.FormFile("file")
		if err != nil {
			io.WriteString(w, "internel server error")
			return
		}
		defer file.Close()
		fileMeta := meta.FileMeta{
			FileSha1: "",
			FileName: header.Filename,
			FileSize: header.Size,
			Location: "/tmp/YunDisk/" + header.Filename,
			UploadAt: time.Now().Format("2006-01-02 15:04:05"),
		}
		newFile, err := os.Create("/YunDisk/" + header.Filename)
		if err != nil {
			io.WriteString(w, "internel server error")
			return
		}
		defer newFile.Close()
		_, err = io.Copy(newFile, file)
		if err != nil {
			io.WriteString(w, "internel server error")
			return
		}
		newFile.Seek(0, 0)
		fileMeta.FileSha1 = util.FileSha1(newFile)
		meta.UpdateFileMeta(fileMeta)
		http.Redirect(w, r, "/file/upload/suc", http.StatusFound)
	}
}

// 上传成功
func UpLoadSucHandler(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Uploaded successfully")
}
