package db

import (
	"filestore-server/db/mysql"
	"fmt"
	"time"
)

type UserFile struct {
	fileSha1 string
	fileSize int64
	fileName string
	uploadAt string
	lastUpload string
}

func OnUserFileUploadFinish (username string, fileSha1 string, fileSize int64, fileName string)bool {
	sql := "INSERT INTO `tbl_user_file`(`user_name`,`file_sha1`,`file_size`,`file_name`,`last_update`) VALUES (?,?,?,?,?)"
	stmt,err := mysql.DbConn().Prepare(sql)
	if err != nil {
		fmt.Printf(" OnUserFileUploadFinish Prepare err:%s\n",err)
		return false
	}
	defer stmt.Close()

	_,err = stmt.Exec(username, fileSha1, fileSize, fileName,time.Now())
	if err != nil {
		fmt.Printf(" OnUserFileUploadFinish stmt.Exec err:%s\n",err)
		return false
	}
	return true
}



func QueryUserFileMeta (username string,limit int) {
	sql := "SELECT * FROM `tbl_user_file` WHERE `user_name` = ? LIMIT ? "
	stmt,err := mysql.DbConn().Prepare(sql)
	if err != nil {
		fmt.Printf(" QueryUserFileMeta Prepare err:%s\n",err)
		return nil
	}
	defer stmt.Close()
}