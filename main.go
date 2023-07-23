package main

import (
	"bytes"
	"compress/gzip"
	"encoding/binary"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"openai-proxy/database"
	"openai-proxy/jwt"
	"os"
	"strings"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	OpenApiKey := os.Getenv("OPEN_API_KEY")
	log.Println("OPEN_API_KEY: ", OpenApiKey)
	SecretKey := os.Getenv("SECRET_KEY")
	jwt.Secret = []byte(SecretKey)

	// 初始化数据库
	database.Init()

	defer database.Close()

	// 创建反向代理的目标URL
	targetURL, err := url.Parse("https://api.openai.com")
	if err != nil {
		log.Fatal(err)
	}

	// 创建反向代理的处理程序
	proxy := httputil.NewSingleHostReverseProxy(targetURL)

	// 解决在容器中报x509:certificate signed by unknown authority
	//proxy.Transport = &http.Transport{TLSClientConfig: &tls.Config{
	//	// InsecureSkipVerify 控制client端是否验证服务器的证书链和主机名，
	//	// 如果InsecureSkipVerify为true, crypto/tls接受服务器提供的任务和证书中的任何主机名
	//	// 在这种模式下，tls容易收到中间机器的攻击，除非是用自定义认证，建议只用于测试或者与VerifyConnection或VerifyPeerCertificate一起使用
	//	InsecureSkipVerify: true,
	//}}

	proxy.Director = func(req *http.Request) {
		//originalDirector(req)
		req.URL.Scheme = targetURL.Scheme
		req.URL.Host = targetURL.Host
		req.Host = targetURL.Host
		req.Header.Set("Authorization", "Bearer "+OpenApiKey)
	}
	// 设置自定义的修改请求的处理函数
	//proxy.ModifyResponse = func(resp *http.Response) error {
	//	// 在此处可以对响应进行修改或记录
	//	// 这里是一个示例，记录响应数据大小
	//	log.Printf("Response size: %d bytes\n", resp.ContentLength)
	//	return nil
	//}
	proxy.ErrorHandler = func(writer http.ResponseWriter, request *http.Request, e error) {
		capture := writer.(*responseCapture)
		log.Printf("%s Proxy %s Error: %s \n", capture.name, capture.path, e)
	}

	// 添加中间件，用于记录每个请求的响应数据大小
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Error reading request body", http.StatusInternalServerError)
			return
		}
		// 日志记录请求数据
		log.Printf("Request: %s %s, Headers: %s,Body: %s \n", r.Method, r.URL, r.Header, body)
		authorization := r.Header.Get("Authorization")
		// 自定义的令牌校验
		if authorization != "" {
			// authorization 按“Bearer ”切分
			stringSlice := strings.Split(authorization, " ")
			if len(stringSlice) >= 1 {
				name := jwt.ValidJwt(stringSlice[1])
				if name != nil {

					// 将处理后的请求主体写入新的请求
					r.Body = io.NopCloser(bytes.NewBuffer(body))
					r.ContentLength = int64(len(body))

					// 创建一个响应捕获器
					captureResponse := &responseCapture{ResponseWriter: w, reqCount: r.ContentLength, name: *name, path: r.URL.Path}
					// 代理转发请求
					proxy.ServeHTTP(captureResponse, r)
					return
				}
			}
		}

		w.WriteHeader(401)
	})

	// 创建HTTP服务器
	server := &http.Server{
		Addr:    ":8080",
		Handler: handler,
	}

	// 启动服务器
	log.Println("Proxy server listening on :8080")
	log.Fatal(server.ListenAndServe())
}

// 自定义响应捕获器，用于记录响应数据大小
type responseCapture struct {
	http.ResponseWriter
	path     string
	reqCount int64
	name     string
}

func (r *responseCapture) Write(b []byte) (int, error) {
	contentEncoding := r.ResponseWriter.Header().Get("Content-Encoding")
	if contentEncoding == "gzip" {
		parseGzip, _ := ParseGzip(b)
		log.Printf("%s Request %s, request size %d bytes, response size: %d bytes, Body: %s\n", r.name, r.path, r.reqCount, len(b), string(parseGzip))
	} else {
		log.Printf("%s Request %s, request size %d bytes, response size: %d bytes\n", r.name, r.path, r.reqCount, len(b))
	}
	// 记录响应数据大小

	database.Insert(r.name, r.path, r.reqCount, int64(len(b)))
	r.ResponseWriter.Header()
	return r.ResponseWriter.Write(b)
}

func ParseGzip(data []byte) ([]byte, error) {
	b := new(bytes.Buffer)
	binary.Write(b, binary.LittleEndian, data)
	r, err := gzip.NewReader(b)
	if err != nil {
		log.Printf("[ParseGzip] NewReader error: %v, maybe data is ungzip \n", err)
		return nil, err
	} else {
		defer r.Close()
		undatas, err := ioutil.ReadAll(r)
		if err != nil {
			log.Printf("[ParseGzip]  ioutil.ReadAll error: %v \n", err)
			return nil, err
		}
		return undatas, nil
	}
}
