package db

import (
	"database/sql"
	"filestore-server/db/mysql"
	"fmt"
)

//更新用户文件表
func OnFileUploadFinished(fileSha1 string, fileName string, fileSize int64, fileAddr string) bool {
	sql := "INSERT INTO `tbl_file`(`file_sha1`,`file_name`,`file_size`,`file_addr`,`status`) VALUES (?,?,?,?,1)"
	stmt,err := mysql.DbConn().Prepare(sql)
	if err != nil {
		fmt.Printf(" db file add Prepare err:%s\n",err)
		return false
	}

	defer stmt.Close()
	ret,err := stmt.Exec(fileSha1, fileName, fileSize, fileAddr)
	if err != nil {
		fmt.Printf(" db file add Exec err:%s\n",err)
		return false
	}

	if rf,err := ret.RowsAffected();err == nil {
		if rf <= 0 {
			fmt.Printf(" db file add rf <= 0, hash:%s\n", fileSha1)
		}
		return true;
	}else{
		fmt.Printf(" db file add RowsAffected err:%s\n",err)
	}

	return false
}

type TableFile struct {
	FileSha1 string
	FileName sql.NullString
	FileSize sql.NullInt64
	FileAddr sql.NullString
}

//获取元信息
func GetFileMeta (fileSha1 string)(*TableFile,error) {
	sql := "SELECT `file_sha1`,`file_name`,`file_size`,`file_addr` FROM `tbl_file` WHERE `file_sha1` = ? AND `status` = 1 LIMIT 1"
	stmt,err := mysql.DbConn().Prepare(sql)
	if err != nil {
		fmt.Printf(" GetFileMeta Prepare err :%s \n",err)
		return nil,err
	}
	tableFile := TableFile{}
	defer stmt.Close()
	err = stmt.QueryRow(fileSha1).Scan(&tableFile.FileSha1,&tableFile.FileName,&tableFile.FileSize,&tableFile.FileAddr)
	if err != nil {
		fmt.Printf(" GetFileMeta QueryRow Scan err :%s \n",err)
		return nil,err
	}
	return &tableFile,nil
}