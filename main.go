package main

import (
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"openai-proxy/src"
	"os"
	"strings"
)

func main() {
	jwt.CreateJwt("TestUser")
	//println(tokenString)
	//ValidJwt(tokenString + "1")
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	OpenApiKey := os.Getenv("OPEN_API_KEY")
	log.Println("OPEN_API_KEY: ", OpenApiKey)

	// 创建反向代理的目标URL
	targetURL, err := url.Parse("https://api.openai.com")
	if err != nil {
		log.Fatal(err)
	}

	// 创建反向代理的处理程序
	proxy := httputil.NewSingleHostReverseProxy(targetURL)
	//originalDirector := proxy.Director
	proxy.Director = func(req *http.Request) {
		//originalDirector(req)
		req.URL.Scheme = targetURL.Scheme
		req.URL.Host = targetURL.Host
		req.Host = targetURL.Host
		req.Header.Set("Authorization", "Bearer "+OpenApiKey)
	}
	// 设置自定义的修改请求的处理函数
	proxy.ModifyResponse = func(resp *http.Response) error {
		// 在此处可以对响应进行修改或记录
		// 这里是一个示例，记录响应数据大小
		log.Printf("Response size: %d bytes\n", resp.ContentLength)
		return nil
	}

	// 添加中间件，用于记录每个请求的响应数据大小
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authorization := r.Header.Get("Authorization")
		// 自定义的令牌校验
		if authorization != "" {
			// authorization 按“Bearer ”切分
			stringSlice := strings.Split(authorization, " ")
			if len(stringSlice) >= 1 && jwt.ValidJwt(stringSlice[1]) {
				// 创建一个响应捕获器
				captureResponse := &responseCapture{ResponseWriter: w}

				// 代理转发请求
				proxy.ServeHTTP(captureResponse, r)
				return
			}
		}

		w.WriteHeader(401)
	})

	// 启动代理服务器
	log.Println("Starting reverse proxy on :8080")
	err = http.ListenAndServe(":8080", handler)
	if err != nil {
		log.Fatal(err)
	}
}

// 自定义响应捕获器，用于记录响应数据大小
type responseCapture struct {
	http.ResponseWriter
}

func (r *responseCapture) Write(b []byte) (int, error) {
	// 记录响应数据大小
	log.Printf("Captured response size: %d bytes\n", len(b))
	return r.ResponseWriter.Write(b)
}
