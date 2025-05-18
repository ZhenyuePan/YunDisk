package handler

import (
	"YunDisk/meta"
	"YunDisk/util"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

// index页面
func UpLoadView(w http.ResponseWriter, r *http.Request) {
	data, err := os.ReadFile("static/view/index.html")
	if err != nil {
		log.Printf("Create error! err: %v", err)
		return
	}

	_, _ = io.WriteString(w, string(data))
}

// 上传文件
func UpLoadHandler(w http.ResponseWriter, r *http.Request) {
	filepath := "/tmp/YunDisk/"
	exist, errA := util.PathExists(filepath)
	if errA != nil {
		log.Printf("PathExists error! err: %v", errA)
		return
	}
	if !exist {
		errB := os.MkdirAll(filepath, os.ModePerm)
		if errB != nil {
			log.Printf("MkdirAll error! err: %v", errB)
			return
		}
	}
	fmt.Println(r.Header)
	file, header, err := r.FormFile("file")
	if err != nil {
		log.Printf("FormFile error! err: %v", err)
		return
	}

	defer file.Close()

	fileMeta := meta.FileMeta{
		FileSha1: "",
		FileName: header.Filename,
		FileSize: header.Size,
		Location: filepath + header.Filename,
		UploadAt: time.Now().Format("2006-01-02 15:04:05"),
	}

	newFile, err := os.Create(filepath + header.Filename)
	if err != nil {
		log.Printf("Create error! err: %v", err)
		return
	}
	defer newFile.Close()
	if _, err = io.Copy(newFile, file); err != nil {
		log.Printf("io Copy error! err: %v", err)
		return
	}

	_, _ = newFile.Seek(0, 0)
	fileMeta.FileSha1 = util.FileSha1(newFile)
	meta.UpdateFileMeta(fileMeta)
	http.Redirect(w, r, "/file/upload/suc", http.StatusFound)
}

// 上传成功
func UpLoadSucHandler(w http.ResponseWriter, r *http.Request) {
	// 设置中文编码头
	w.Header().Set("Content-Type", "text/plain; charset=utf-8md")
	_, _ = io.WriteString(w, "文件传输成功！") // 修改输出内容
}

// 获取文件元信息
func GetFileMetaHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	filehash := r.Form["filehash"][0]
	fMeta := meta.GetFileMeta(filehash)
	data, err := json.Marshal(fMeta)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(data)
}
