package db

import (
	"filestore-server/db/mysql"
	"fmt"
	"time"
)

type UserFile struct {
	FileSha1 string
	FileSize int64
	FileName string
	UploadAt string
	LastUpdate string
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



func QueryUserFileMeta (username string,limit int) ([]UserFile,error) {
	sql := "SELECT `file_sha1`,`file_size`,`file_name`,`upload_at`,`last_update` FROM `tbl_user_file` WHERE `user_name` = ? LIMIT ? "
	stmt,err := mysql.DbConn().Prepare(sql)
	if err != nil {
		fmt.Printf(" QueryUserFileMeta Prepare err:%s\n",err)
		return nil,err
	}
	defer stmt.Close()

	rows,err := stmt.Query(username,limit)
	if err != nil {
		fmt.Printf(" QueryUserFileMeta stmt.Query err:%s\n",err)
		return nil,err
	}

	var userfiles []UserFile
	for rows.Next() {
		ufile := UserFile{}
		err = rows.Scan(&ufile.FileSha1,&ufile.FileSize,&ufile.FileName,&ufile.UploadAt,&ufile.LastUpdate)
		if err != nil {
			fmt.Printf(" QueryUserFileMeta rows.Scan err:%s\n",err)
			return userfiles,err
		}
		userfiles = append(userfiles, ufile)
	}
	return userfiles,nil
}