package handler

import (
	"YunDisk/meta"
	"YunDisk/util"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"time"
)

func UpLoadView(w http.ResponseWriter, r *http.Request) {
	data, err := os.ReadFile("static/view/index.html")
	if err != nil {
		log.Printf("Create error! err: %v", err)
		return
	}

	_, _ = io.WriteString(w, string(data))
}

func UpLoadHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.Header)
	file, header, err := r.FormFile("file")
	if err != nil {
		log.Printf("FormFile error! err: %v", err)
		return
	}

	defer func(file multipart.File) {
		_ = file.Close()
	}(file)

	fileMeta := meta.FileMeta{
		FileSha1: "",
		FileName: header.Filename,
		FileSize: header.Size,
		Location: "/tmp/YunDisk/" + header.Filename,
		UploadAt: time.Now().Format("2006-01-02 15:04:05"),
	}

	newFile, err := os.Create("/tmp/YunDisk/" + header.Filename)
	if err != nil {
		log.Printf("Create error! err: %v", err)
		return
	}

	defer func(newFile *os.File) {
		_ = newFile.Close()
	}(newFile)

	if _, err = io.Copy(newFile, file); err != nil {
		log.Printf("io Copy error! err: %v", err)

		return
	}

	_, _ = newFile.Seek(0, 0)
	fileMeta.FileSha1 = util.FileSha1(newFile)
	meta.UpdateFileMeta(fileMeta)

	http.Redirect(w, r, "/file/upload/suc", http.StatusFound)
}

func UpLoadSucHandler(w http.ResponseWriter, r *http.Request) {
	_, _ = io.WriteString(w, "Uploaded successfully")
}
