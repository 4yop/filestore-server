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

	http.HandleFunc("/file/downloadurl",handler.DownloadHandler)

	http.HandleFunc("/user/signup",handler.SignUpHandler)

	http.HandleFunc("/user/signin",handler.SignInHandler)

	http.HandleFunc("/static/view/signin.html",handler.SignInHandler)
	http.HandleFunc("/static/view/home.html",handler.HomeHandler)

	http.HandleFunc("/user/info",handler.HTTPInterceptor(handler.UserInfoHandle))
	http.HandleFunc("/file/query",handler.FileQueryHandler)


	// 秒传接口
	http.HandleFunc("/file/fastupload", handler.HTTPInterceptor(
		handler.TryFastUploadHandler))



	// 分块上传接口
	http.HandleFunc("/file/mpupload/init",
		handler.HTTPInterceptor(handler.InitMulitPartUploadHandler))
	http.HandleFunc("/file/mpupload/uppart",
		handler.HTTPInterceptor(handler.UploadPartHandler))
	http.HandleFunc("/file/mpupload/complete",
		handler.HTTPInterceptor(handler.CompleteUoloadHandler))

	http.HandleFunc("/file/mpupload/cancel",
		handler.HTTPInterceptor(handler.CancelUploadHandler))

	http.HandleFunc("/file/mpupload/status",
		handler.HTTPInterceptor(handler.MulitPartUploadHandler))


	err := http.ListenAndServe(":8080",nil)
	if err != nil {
		fmt.Printf("server error:%s \n",err)
	}
}

