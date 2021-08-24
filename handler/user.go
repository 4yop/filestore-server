package handler

import (
	"filestore-server/db"
	"filestore-server/util"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"time"
)

const (
	pwd_salt = "^#$%#8045"
)

func SignUpHandler(w http.ResponseWriter,r *http.Request)  {
	if r.Method == http.MethodGet{
		data,err := ioutil.ReadFile("./static/view/signup.html")
		if err != nil {
			fmt.Println(err)
			io.WriteString(w,"inter server error")
			return
		}
		io.WriteString(w,string(data))
	}else if r.Method == http.MethodPost {
		r.ParseForm()
		username := r.Form.Get("username")
		password := r.Form.Get("password")
		if len(username) < 5 {
			w.Write([]byte("username 太短了"))
			return
		}
		if len(password) < 5 {
			w.Write([]byte("password 太短了"))
			return
		}
		enc_pwd := util.Sha1([]byte(password+pwd_salt))
		if db.UserSignUp(username,enc_pwd) {
			w.Write([]byte("SUCCESS"))
			return
		}else {
			w.Write([]byte("error"))
			return
		}
	}

}


func SignInHandler (w http.ResponseWriter,r *http.Request) {
	if r.Method == http.MethodGet{
		data,err := ioutil.ReadFile("./static/view/signin.html")
		if err != nil {
			fmt.Println(err)
			io.WriteString(w,"inter server error")
			return
		}
		io.WriteString(w,string(data))
	}else if r.Method == http.MethodPost {
		r.ParseForm()
		username := r.Form.Get("username")
		password := r.Form.Get("password")
		if len(username) < 5 {
			w.Write([]byte("username 太短了"))
			return
		}
		if len(password) < 5 {
			w.Write([]byte("password 太短了"))
			return
		}
		enc_pwd := util.Sha1([]byte(password+pwd_salt))
		if !db.UserSignIn(username,enc_pwd) {
			w.Write([]byte("error"))
			return
		}
		token := GenToken(username)
		if !db.UpdateToken(username,token) {
			w.Write([]byte("error"))
			return
		}




		resp := util.RespMsg{
			Code: 0,
			Msg:  "OK",
			Data: map[string]string{
				"Token":token,
				"Username":username,
				"Location" : "http://"+r.Host+"/static/view/home.html",
			},
		}

		w.Write(resp.JSONBytes())
	}
}

func GenToken (username string) string {
	ts := fmt.Sprintf("%x",time.Now().Unix())
	tokenPrefix := util.MD5([]byte(username+ts+"_token_salt"))
	return tokenPrefix + ts[:8]
}

func UserInfoHandle (w http.ResponseWriter,r *http.Request) {
	r.ParseForm()
	username := r.Form.Get("username")
	//token := r.Form.Get("token")
	//if  db.CheckToken(username,token) == false {
	//
	//	w.Write([]byte("error"))
	//	return
	//}
	user,err := db.GetByUsername(username)
	if err != nil {
		w.Write([]byte("error"))
		return
	}
	resp := util.RespMsg{
		Code: 0,
		Msg:  "OK",
		Data: user,
	}
	w.Write(resp.JSONBytes())
}