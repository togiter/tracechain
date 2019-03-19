package webServer

import (
	"fmt"
	"github.com/tracechain/web-service/controller"
	"net/http"
)

func WebStart(app *controller.Application) {
	fmt.Println(http.Dir("web-service/static"))
	fs := http.FileServer(http.Dir("web-service/static"))
	http.Handle("/html/", fs)
	http.Handle("/js/", fs)
	http.Handle("/css/", fs)
	http.HandleFunc("/issueProduct", app.IssueProduct)
	http.HandleFunc("/queryProducts", app.QueryProducts)
	http.HandleFunc("/queryProductNo", app.QueryProductNo)
	http.HandleFunc("/queryProductRange", app.QueryProductRange)
	http.HandleFunc("/transferProduct", app.TransferProduct)
	http.HandleFunc("/alterProductPrice", app.AlterProductPrice)

	http.HandleFunc("/issue.html", app.ShowIssueHtml)
	fmt.Println("启动服务器监听,监听端口:3000")
	err := http.ListenAndServe(":3000", nil)
	if err != nil {
		fmt.Println("启动web服务失败！")
	}
}
