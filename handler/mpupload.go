package handler

import "net/http"

//初始 分块上传
func InitMulitPartUploadHandler (w http.ResponseWriter, r *http.Request) {
	//1.解析请求参数
	//2.获取redis conn
	//3.生成分块的初始化信息
	//4.初始化信息写入缓存
	//5.响应初始化信息返回客户端
}
