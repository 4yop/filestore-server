package main

import (
	"filestore-server/handler"
	"fmt"
	"net/http"
)

func main () {

	http.Handle("/static/",
		http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

	http.HandleFunc("/file/upload",handler.UploadHandler)

	http.HandleFunc("/file/upload/success",handler.UploadSucHandler)

	http.HandleFunc("/file/meta",handler.GetFileMetaHandler)
	http.HandleFunc("/file/download",handler.DownloadHandler)

	http.HandleFunc("/user/signup",handler.SignUpHandler)

	http.HandleFunc("/user/signin",handler.SignInHandler)

	http.HandleFunc("/static/view/signin.html",handler.SignInHandler)
	http.HandleFunc("/static/view/home.html",handler.HomeHandler)

	http.HandleFunc("/user/info",handler.UserInfoHandle)

	err := http.ListenAndServe(":8080",nil)
	if err != nil {
		fmt.Printf("server error:%s \n",err)
	}
}

