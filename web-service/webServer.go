package webServer

import (
	"fmt"
	"github.com/tracechain/web-service/controller"
	"net/http"
)

func WebStart(app *controller.Application) {
	fs := http.FileServer(http.Dir("web-service/static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	http.HandleFunc("/issueProduct", app.IssueProduct)
	http.HandleFunc("/queryProductNo", app.QueryProductNo)
	http.HandleFunc("/queryProductRange", app.QueryProductRange)
	http.HandleFunc("/transferProduct", app.TransferProduct)
	http.HandleFunc("/alterProductPrice", app.AlterProductPrice)
	fmt.Println("启动服务器监听,监听端口:9000")
	err := http.ListenAndServe(":9000", nil)
	if err != nil {
		fmt.Println("启动web服务失败！")
	}
}
