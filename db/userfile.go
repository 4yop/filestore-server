package db

import (
	"filestore-server/db/mysql"
	"fmt"
)

func OnUserFileUploadFinish (username string,file_sha1 string,file_size int64,file_name string) {
	sql := "INSERT INTO `tbl_user_file`(`user_name`,`file_sha1`,`file_size`,`file_name`) VALUES (?,?,?,?)"
	stmt,err := mysql.DbConn().Prepare(sql)
	if err != nil {
		fmt.Printf(" OnUserFileUploadFinish Prepare err:%s\n",err)
	}
	defer stmt.Close()



}
