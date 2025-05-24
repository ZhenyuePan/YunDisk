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
	w.Header().Set("Content-Type", "text/plain; charset=utf-8") // 修正编码声明
	w.WriteHeader(http.StatusOK)
	_, _ = io.WriteString(w, "文件上传成功！") // 明确提示成功状态
}

// 获取文件元信息
// 未测试
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

// 下载文件
// 未测试
func DownloadHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	fsha1 := r.Form.Get("filehash")
	fm := meta.GetFileMeta(fsha1)
	f, err := os.Open(fm.Location)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer f.Close()
	data, err := io.ReadAll(f)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("content-disposition", "attachment; filename=\""+fm.FileName+"\"")
	w.Write(data)
}

// 更新文件元信息
// 未测试
func FileMetaUpdateHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	opType := r.Form.Get("op")
	fileSha1 := r.Form.Get("filehash")
	newFileName := r.Form.Get("filename")
	if opType != "0" {
		w.WriteHeader(http.StatusForbidden)
		return
	}
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	curFileMeta := meta.GetFileMeta(fileSha1)
	curFileMeta.FileName = newFileName
	meta.UpdateFileMeta(curFileMeta)
	w.WriteHeader(http.StatusOK)
	data, err := json.Marshal(curFileMeta)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

// 删除文件及元信息
// 未测试
func FileDeleteHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	fileSha1 := r.Form.Get("filehash")
	fMeta := meta.GetFileMeta(fileSha1)
	os.Remove(fMeta.Location)
	meta.RemoveFileMeta(fileSha1)

	w.WriteHeader(http.StatusOK)
}
