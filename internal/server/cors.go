package server

import (
	"net/http"
)

// CORSHandler 是一个中间件，用于处理 CORS 请求
func CORSHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 设置 CORS 头部
		w.Header().Set("Access-Control-Allow-Origin", "*")                            // 允许所有来源
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")          // 允许的 HTTP 方法
		w.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type") // 允许的请求头

		// 如果是预检请求，直接返回
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		// 继续处理请求
		next.ServeHTTP(w, r)
	})
}
