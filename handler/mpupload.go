package handler

import (
	rPool "filestore-server/cache/redis"
	"filestore-server/util"
	"fmt"
	"github.com/garyburd/redigo/redis"
	"math"
	"net/http"
	"os"
	"strconv"
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
	//3.获取文件句柄，存储区块内容
	fd,err := os.Create("./tmp/"+uoloadId+"/"+chunkIndex)
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
	total := 0
	chunkcount := 0
	

	//4.合并， copy 等方式


	//4.改 数据 的状态
	rConn.Do("HDEL","MP_"+uoloadId)
	//5.返回结果
	w.Write(util.NewRespMsg(0,"OK",nil).JSONBytes())
}