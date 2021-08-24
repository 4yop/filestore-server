package db

import (
	"filestore-server/db/mysql"
	"fmt"
)

//用户注册
func UserSignUp (username string,password string) bool{
	sql := "INSERT INTO `tbl_user`(`user_name`,`user_pw`) VALUES (?,?)"
	stmt,err := mysql.DbConn().Prepare(sql)
	if err != nil {
		fmt.Printf(" UserSignUp Prepare  err:%s\n",err)
		return false
	}
	defer stmt.Close()
	ret,err := stmt.Exec(username,password)
	if err != nil {
		fmt.Printf(" UserSignUp stmt Exec  err:%s\n",err)
		return false
	}
	row,err := ret.RowsAffected();
	if err != nil {
		fmt.Printf(" UserSignUp RowsAffected  err:%s\n",err)
		return false
	}
	if row < 1 {
		fmt.Printf(" UserSignUp ret.RowsAffected row < 1  ",)
		return false
	}
	fmt.Printf("UserSignUp success:user_name:%s,user_pw:%s\n",username,password)
	return true
}

//用户登录
func UserSignIn (username string,password string) bool {
	sql := "SELECT * FROM `tbl_user` WHERE `user_name` = ? LIMIT 1"
	stmt,err := mysql.DbConn().Prepare(sql)
	if err != nil {
		fmt.Printf("UserSignIn stmt err:%s\n",err)
		return false
	}

	defer stmt.Close()
	rows,err := stmt.Query(username)
	if err != nil {
		fmt.Printf("UserSignIn Query err:%s\n",err)
		return false
	}
	if rows == nil {
		fmt.Printf("UserSignIn Query rows == nil ",)
		return false
	}

	//pRows := rows.Columns("user_pw")
	pRows := mysql.ParseRows(rows)
	if len(pRows) > 0 && string(pRows[0]["user_pw"].([]byte)) == password {
		return true
	}
	fmt.Printf("密码不对")
	return false
}

func UpdateToken(username string,token string) bool {
	sql := "REPLACE INTO `tbl_token`(`user_name`,`user_token`) VALUES (?,?)"
	stmt,err := mysql.DbConn().Prepare(sql)
	if err != nil {
		fmt.Printf("UpdateToken Prepare err:%s\n",err)
		return false
	}
	defer stmt.Close()
	_,err = stmt.Exec(username,token)
	if err != nil {
		fmt.Printf("UpdateToken Exec err:%s\n",err)
		return false
	}
	//rows,err := ret.RowsAffected()
	//if err != nil {
	//	fmt.Printf("UpdateToken ret.RowsAffected err:%s\n",err)
	//	return false
	//}
	//if rows < 1 {
	//	fmt.Printf("UpdateToken ret.RowsAffected.rows < 1 \n")
	//	return false
	//}
	return true
}

//检查token
func CheckToken(username string,token string) bool {
	sql := "SELECT `user_name`,`user_token` FROM `tbl_token` WHERE `user_name` = ? AND  `user_token` = ? LIMIT 1"
	stmt,err := mysql.DbConn().Prepare(sql)
	if err != nil {
		fmt.Printf("CheckToken Prepare err:%s\n",err)
		return false
	}

	defer stmt.Close()
	rows,err := stmt.Query(username,token)
	if err != nil {
		fmt.Printf("CheckToken Query err:%s\n",err)
		return false
	}
	if rows == nil {
		fmt.Printf("CheckToken Query rows == nil ",)
		return false
	}
	return true
}

type User struct {
	Username string
	Email string
	Phone string
	SignupAt string
	LastActiveAt string
	Status int
}

func GetByUsername(username string) (User,error){
	sql := "SELECT `user_name`,`signup_at` FROM `tbl_user` WHERE `user_name` = ? LIMIT 1 "
	stmt,err := mysql.DbConn().Prepare(sql)
	if err != nil {
		fmt.Printf(" GetByUsername Prepare err:%s \n",err)
		return User{},err
	}
	defer stmt.Close()
	user := User{}
	err = stmt.QueryRow(username).Scan(&user.Username,&user.SignupAt)
	if err != nil {
		fmt.Printf(" GetByUsername QueryRow.Scan err:%s \n",err)
		return User{},err
	}
	return user,err
}
