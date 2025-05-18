package handler

import (
	"YunDisk/meta"
	"YunDisk/util"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

// 处理文件上传
func UpLoadHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		data, err := os.ReadFile("static/view/index.html")
		if err != nil {
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}
		w.Write(data)

	case "POST":
		// 1. 获取上传文件
		file, header, err := r.FormFile("file")
		if err != nil {
			fmt.Printf("FormFile error: %v\n", err)
			http.Error(w, "failed to get uploaded file", http.StatusBadRequest)
			return
		}
		defer file.Close()

		// 2. 准备目标目录
		targetDir := "/tmp/YunDisk"
		if err := os.MkdirAll(targetDir, 0755); err != nil {
			fmt.Printf("MkdirAll error: %v\n", err)
			http.Error(w, "failed to create target directory", http.StatusInternalServerError)
			return
		}

		// 3. 创建目标文件
		targetPath := filepath.Join(targetDir, header.Filename)
		newFile, err := os.Create(targetPath)
		if err != nil {
			fmt.Printf("Create error: %v (Path: %s)\n", err, targetPath)
			http.Error(w, "failed to create file", http.StatusInternalServerError)
			return
		}
		defer newFile.Close()

		// 4. 复制文件内容
		if _, err := io.Copy(newFile, file); err != nil {
			fmt.Printf("Copy error: %v\n", err)
			http.Error(w, "failed to save file content", http.StatusInternalServerError)
			return
		}

		// 5. 计算SHA1
		if _, err := newFile.Seek(0, 0); err != nil {
			fmt.Printf("Seek error: %v\n", err)
			http.Error(w, "failed to read file", http.StatusInternalServerError)
			return
		}
		fileSha1 := util.FileSha1(newFile)

		// 6. 保存元数据
		fileMeta := meta.FileMeta{
			FileSha1: fileSha1,
			FileName: header.Filename,
			FileSize: header.Size,
			Location: targetPath,
			UploadAt: time.Now().Format("2006-01-02 15:04:05"),
		}
		meta.UpdateFileMeta(fileMeta)
		fmt.Printf("File uploaded successfully: %s\n", targetPath)

		http.Redirect(w, r, "/file/upload/suc", http.StatusFound)
	}
}

// 上传成功
func UpLoadSucHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Uploaded successfully"))
}
