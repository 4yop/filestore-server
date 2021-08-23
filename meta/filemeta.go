package meta

import "filestore-server/db"

// FileMeta : 文件元信息结构
type FileMeta struct {
	FileSha1 string
	FileName string
	FileSize int64
	Location string
	UploadAt string
}

var fileMetas map[string]FileMeta

func init() {
	fileMetas = make(map[string]FileMeta)
}

// UpdateFileMeta : 新增/更新文件元信息
func UploadFileMeta (fMeta FileMeta) {
	fileMetas[fMeta.FileSha1] = fMeta
}

// GetFileMeta : 通过sha1值获取文件的元信息对象
func GetFileMeta(fileSha1 string) FileMeta {
	return fileMetas[fileSha1]
}

// GetFileMeta : 通过sha1值获取文件的元信息对象
func GetFileMetaDb(fileSha1 string) (FileMeta,error) {
	tableFile,err := db.GetFileMeta(fileSha1)
	if err != nil {
		return FileMeta{},err
	}
	return FileMeta{
		FileSha1 : tableFile.FileSha1,
		FileName : tableFile.FileName.String,
		FileSize : tableFile.FileSize.Int64,
		Location : tableFile.FileAddr.String,
	},nil
}

//元数据保存到数据库
func UpdateFileMetaDb(fMeta FileMeta) bool {
	return db.OnFileUploadFinished(
		fMeta.FileSha1,
		fMeta.FileName,
		fMeta.FileSize,
		fMeta.Location)
}