package main

import (
	"flag"
	"fmt"
	"net/http"

	"github.com/gythialy/jrebel/constant"
	"github.com/gythialy/jrebel/handler"
)

func main() {
	var host string
	var port string
	flag.StringVar(&port, "p", "9000", "端口,默认为9000")
	flag.StringVar(&host, "h", "0.0.0.0", "绑定host,默认为0.0.0.0")
	flag.Parse()

	leaseHandler := handler.NewHandler()

	// 处理根路径和UUID路径
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		if path == "/" {
			w.WriteHeader(http.StatusOK)
			w.Header().Add("Content-Type", "text/html; charset=utf-8")
			_, _ = w.Write([]byte("<h1>Hello, jrebel!</h1>"))
		} else if len(path) > 1 {
			uuid := path[1:]
			if isValidUUID(uuid) {
				leaseHandler.Leases(w, r)
			} else {
				http.NotFound(w, r)
			}
		} else {
			http.NotFound(w, r)
		}
	})

	http.HandleFunc("/uuid", handler.UUID)
	http.HandleFunc("/jrebel/leases", leaseHandler.Leases)
	http.HandleFunc("/jrebel/leases/1", leaseHandler.Leases1)
	http.HandleFunc("/agent/leases", leaseHandler.Leases)
	http.HandleFunc("/agent/leases/1", leaseHandler.Leases1)
	http.HandleFunc("/jrebel/validate-connection", leaseHandler.ValidateConnection)
	http.HandleFunc("/rpc/ping.action", handler.PingHandler)
	http.HandleFunc("/rpc/obtainTicket.action", handler.ObtainTicketHandler)
	http.HandleFunc("/rpc/releaseTicket.action", handler.ReleaseTicketHandler)

	fmt.Printf(`
	启动成功 端口号: %s
GET /uuid 生成随机串
http://%s:%s/{uuid} 放入jrebel激活地址栏
Vesion: %s
BuildTime:%s
`, port, host, port, constant.Version, constant.BuildTime)

	err := http.ListenAndServe(host+":"+port, nil)
	if err != nil {
		fmt.Printf("http.ListenAndServe() 函数执行错误,错误为:%v\n", err)
		return
	}
}

// isValidUUID UUID验证函数
func isValidUUID(uuid string) bool {
	if len(uuid) != 36 {
		return false
	}
	for i, c := range uuid {
		if i == 8 || i == 13 || i == 18 || i == 23 {
			if c != '-' {
				return false
			}
		} else {
			if (c < '0' || c > '9') && (c < 'a' || c > 'f') && (c < 'A' || c > 'F') {
				return false
			}
		}
	}
	return true
}
