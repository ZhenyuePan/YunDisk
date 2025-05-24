package meta

// 文件元信息结构
type FileMeta struct {
	FileSha1 string //哈希值
	FileName string //文件名
	FileSize int64  //文件大小
	Location string //文件存储路径
	UploadAt string //时间戳
}

var fileMetas map[string]FileMeta

func init() {
	fileMetas = make(map[string]FileMeta)
}

// 更新/新增文件元信息
func UpdateFileMeta(fmeta FileMeta) {
	fileMetas[fmeta.FileSha1] = fmeta
}

// 通过sha1值获取文件的元信息对象
func GetFileMeta(fileSha1 string) FileMeta {
	return fileMetas[fileSha1]
}

// 通过sha1删除文件元信息
func RemoveFileMeta(fileSha1 string) {
	delete(fileMetas, fileSha1)
}
