package handler

import (
	"encoding/json"
	"filestore-server/db"
	"filestore-server/meta"
	"filestore-server/util"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

//上传
func UploadHandler(w http.ResponseWriter,r *http.Request)  {
	if r.Method == http.MethodGet {
		fmt.Println(r.Method)
		//查询

		data,err := ioutil.ReadFile("./static/view/index.html")
		if err != nil {
			fmt.Println(err)
			io.WriteString(w,"inter server error")
			return
		}
		io.WriteString(w,string(data))
	}else if r.Method == http.MethodPost {
		//上传 等
		file,head,err := r.FormFile("file")
		if err != nil {
			fmt.Printf("上传失败,err:%s \n",err)
			return
		}
		defer file.Close()

		//util.FileSha1(file)
		fMeta := meta.FileMeta{
				FileSha1: "",
				FileName: head.Filename,
				Location: "./tmp/"+head.Filename,
				UploadAt: time.Now().Format("2021-8-20 16:32:01"),
			}

		newFile,err := os.Create(fMeta.Location)
		if err != nil {
			fmt.Printf("创建失败,err:%s \n",err)
			return
		}
		defer newFile.Close()

		fMeta.FileSize,err = io.Copy(newFile,file)
		if err != nil {
			fmt.Printf("保存失败,err:%s \n",err)
			return
		}

		newFile.Seek(0,0)
		fMeta.FileSha1 = util.FileSha1(newFile)
		fmt.Println(fMeta)
		//meta.UploadFileMeta(fMeta)
		meta.UpdateFileMetaDb(fMeta)
		r.ParseForm()
		username := r.Form.Get("username")
		suc := db.OnUserFileUploadFinish(username,fMeta.FileSha1,fMeta.FileSize,fMeta.FileName)
		if suc {
			http.Redirect(w,r,"/file/upload/success",http.StatusFound)
		}else{
			w.Write([]byte("upload failed"))
		}


	}
}

//获取文件元信息
func GetFileMetaHandler(w http.ResponseWriter,r *http.Request) {
	r.ParseForm()

	filehash := r.Form.Get("filehash")

	//fMeta := meta.GetFileMeta(filehash)

	fMeta,err := meta.GetFileMetaDb(filehash)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return;
	}

	data,err := json.Marshal(fMeta)
	if err != nil {
		fmt.Println("json 转换失败")
		w.WriteHeader(http.StatusInternalServerError)
		return;
	}
	w.Write(data)
}



//成功返回信息
func UploadSucHandler (w http.ResponseWriter,r *http.Request) {
	io.WriteString(w,"Upload finish!")
}


func DownloadHandler (w http.ResponseWriter,r *http.Request)  {
	r.ParseForm()

	filehash := r.Form.Get("filehash")
	fMeta := meta.GetFileMeta(filehash)


	file,err := os.Open(fMeta.Location)
	if err != nil {
		fmt.Println("文件打开失败:"+fMeta.Location)
		w.WriteHeader(http.StatusInternalServerError)
		return;
	}
	defer file.Close()

	data,err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Println("文件读取失败:"+fMeta.Location)
		w.WriteHeader(http.StatusInternalServerError)
		return;
	}

	w.Header().Set("Content-Type","application/octet-stream")
	w.Header().Set("content-disposition","attachment;filename=\""+fMeta.FileName+"\"")
	w.Write(data)
}

func HomeHandler(w http.ResponseWriter,r *http.Request) {
	data,err := ioutil.ReadFile("./static/view/home.html")
	if err != nil {
		fmt.Printf(" home handler ioutil.readfile err:%s\n",err)
	}
	w.Write(data)
}