package handler

import (
	rPool "filestore-server/cache/redis"
	"filestore-server/db"
	"filestore-server/util"
	"fmt"
	"github.com/garyburd/redigo/redis"
	"math"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"
	"time"
)
type MulitPartUploadInfo struct {
	FileHash string
	FileSize int
	UploadID string
	ChunkSize int
	ChunkCount int
}

//初始 分块上传
func InitMulitPartUploadHandler (w http.ResponseWriter, r *http.Request) {
	//1.解析请求参数
	r.ParseForm()
	username := r.Form.Get("username")
	filehash := r.Form.Get("filehash")
	filesize,err := strconv.Atoi(r.Form.Get("filesize"))
	if err != nil {
		w.Write(util.NewRespMsg(-1, "parame invalid", nil).JSONBytes())
		return
	}
	//2.获取redis conn
	rConn := rPool.RedisPool().Get()
	defer rConn.Close()

	//3.生成分块的初始化信息
	upinfo := MulitPartUploadInfo{
		FileHash:   filehash,
		FileSize:   filesize,
		UploadID:   username + fmt.Sprintf("%s",time.Now().UnixNano()),
		ChunkSize:  5*1024*1024,
		ChunkCount: int(math.Ceil(float64(filesize)/ 5 * 1024 * 1024)),
	}
	//4.初始化信息写入缓存
	rConn.Do("HSET","MP_"+upinfo.UploadID,"filehash",upinfo.FileHash)
	rConn.Do("HSET","MP_"+upinfo.UploadID,"filesize",upinfo.FileSize)
	rConn.Do("HSET","MP_"+upinfo.UploadID,"chunksize",upinfo.ChunkSize)
	rConn.Do("HSET","MP_"+upinfo.UploadID,"chunkcount",upinfo.ChunkCount)
	//5.响应初始化信息返回客户端
	w.Write(util.NewRespMsg(0,"ok",upinfo).JSONBytes())
}

//上传文件块
func UploadPartHandler(w http.ResponseWriter, r *http.Request)  {
	//1.获取参数
	r.ParseForm()
	//username := r.Form.Get("username")
	uoloadId := r.Form.Get("uoloadId")
	chunkIndex := r.Form.Get("chunkIndex")
	//2.reids conn
	rConn := rPool.RedisPool().Get()
	defer rConn.Close()
	fpath := "./tmp/"+uoloadId+"/"+chunkIndex
	os.MkdirAll(path.Dir(fpath),0744)
	//3.获取文件句柄，存储区块内容
	fd,err := os.Create(fpath)
	if err != nil {
		fmt.Printf("UploadPartHandler os.Create:%s,err:%s\n","./tmp/"+uoloadId+"/"+chunkIndex,err)
		w.Write(util.NewRespMsg(-1,err.Error(),nil).JSONBytes())
		return;
	}
	defer fd.Close()

	buf := make([]byte,1024*1024)
	for  {
		n,err := r.Body.Read(buf)
		if err != nil {
			fmt.Printf("UploadPartHandler r.Body.Read ,err:%s\n",err)
			//w.Write(util.NewRespMsg(-1,err.Error(),nil).JSONBytes())
			break
		}
		n,err = fd.Write(buf[:n])
		if err != nil {
			fmt.Printf("UploadPartHandler fd.Write ,err:%s\n",err)
			//w.Write(util.NewRespMsg(-1,err.Error(),nil).JSONBytes())
			break
		}
	}

	//4.更新redis的状态
	rConn.Do("HSET","MP_"+uoloadId,"chkidx_"+chunkIndex,1)
	//5.返回处理结果
	w.Write(util.NewRespMsg(0,"OK",nil).JSONBytes())
}

//通知上传合并
func CompleteUoloadHandler(w http.ResponseWriter, r *http.Request) {
	//1.获取参数
	r.ParseForm()
	filehash := r.Form.Get("filehash")
	uoloadId := r.Form.Get("uoloadId")
	username := r.Form.Get("username")
	filesize,_ := strconv.Atoi(r.Form.Get("filesize"))
	filename := r.Form.Get("filename")
	//2.redis conn
	rConn := rPool.RedisPool().Get()

	//3.判断每一块是否都上传
	data,err := redis.Values(rConn.Do("HGETALL","MP_"+uoloadId))
	if err != nil {
		fmt.Printf("CompleteUoloadHandler redis.Values HGETALL MP_%s ,err:%s\n",uoloadId,err)
		w.Write(util.NewRespMsg(-1,err.Error(),nil).JSONBytes())
		return
	}
	totalCount := 0
	chunkcount := 0
	for i := 0; i < len(data); i++ {
		k := string(data[i].([]byte))
		v := string(data[i+1].([]byte))
		if k == "chunkcount" {
			totalCount,_ = strconv.Atoi(v)
		}else if strings.HasPrefix(k,"chkidx_") && v == "1" {
			chunkcount += 1
		}
	}
	if totalCount != chunkcount {
		w.Write(util.NewRespMsg(-2,"未传完所有块",nil).JSONBytes())
	}
	//4.合并， copy 等方式
	tmpName := "./tmp/"+uoloadId+"-"+filehash+".tmp"
	fd,err := os.Create(tmpName)
	if err != nil {
		w.Write(util.NewRespMsg(-2,err.Error(),nil).JSONBytes())
		return;
	}
	defer fd.Close()
	tmp := make([]byte,1024*1024)
	for i := 0; i < totalCount; i++ {
		chunkfile,_ := os.Open("./tmp/"+uoloadId+"/"+strconv.Itoa(i))
		defer chunkfile.Close()
		for  {
			n,err := chunkfile.Read(tmp)
			if err != nil  {
				break
			}
			fd.Write(tmp[:n])
		}
	}
	err = os.Rename(tmpName,"./tmp/"+filename)
	if err != nil {
		w.Write(util.NewRespMsg(-1,err.Error(),nil).JSONBytes())
		return ;
	}

	//4.改 数据 的状态
	db.OnFileUploadFinished(filehash,filename,int64(filesize),"fileaddr")
	db.OnUserFileUploadFinish(username,filehash,int64(filesize),filename)
	rConn.Do("HDEL","MP_"+uoloadId)
	//5.返回结果
	w.Write(util.NewRespMsg(0,"OK",nil).JSONBytes())
}

//取消上传
func CancelUploadHandler (w http.ResponseWriter,r *http.Request) {
	//1.删除分块
	//2.删除redis状态
	//3.更新mysql
	//4.返回结果
}

//上传状态
func MulitPartUploadHandler (w http.ResponseWriter,r *http.Request) {
	//检查分块是否有效
	//获取分块初始化信息
	//获取已上传的分块信息
}