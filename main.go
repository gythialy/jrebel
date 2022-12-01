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
	//1.注册处理器函数
	http.HandleFunc("/uuid", handler.UUID)
	http.HandleFunc("/jrebel/leases", handler.Leases)
	http.HandleFunc("/jrebel/leases/1", handler.Leases1)
	http.HandleFunc("/jrebel/validate-connection", handler.ValidateConnection)
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
