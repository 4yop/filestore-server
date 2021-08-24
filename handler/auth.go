package handler

import "net/http"

//http 请求拦截器
func HTTPInterceptor(f func(w http.ResponseWriter,r *http.Request)) (func(w http.ResponseWriter,r *http.Request)) {
	
	return f
}
