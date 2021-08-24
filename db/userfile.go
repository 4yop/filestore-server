package db

import (
	"filestore-server/db/mysql"
	"fmt"
)

func OnUserFileUploadFinish (username string, fileSha1 string, fileSize int64, fileName string)bool {
	sql := "INSERT INTO `tbl_user_file`(`user_name`,`file_sha1`,`file_size`,`file_name`) VALUES (?,?,?,?)"
	stmt,err := mysql.DbConn().Prepare(sql)
	if err != nil {
		fmt.Printf(" OnUserFileUploadFinish Prepare err:%s\n",err)
		return false
	}
	defer stmt.Close()

	_,err = stmt.Exec(username, fileSha1, fileSize, fileName)
	if err != nil {
		fmt.Printf(" OnUserFileUploadFinish stmt.Exec err:%s\n",err)
		return false
	}
	return true
}
