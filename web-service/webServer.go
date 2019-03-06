package webServer

import (
	"fmt"
	"genealogy/web-service/controller"
	"net/http"
)

func WebStart(app *controller.Application) {
	fs := http.FileServer(http.Dir("web-service/static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	http.HandleFunc("/", app.IndexView)
	http.HandleFunc("/index.html", app.IndexView)
	http.HandleFunc("/add.html", app.AddView)
	http.HandleFunc("/query.html", app.QueryView)
	http.HandleFunc("/addMember", app.AddMember)
	http.HandleFunc("/queryMember", app.QueryMember)
	fmt.Println("启动服务器监听,监听端口:9000")
	err := http.ListenAndServe(":9000", nil)
	if err != nil {
		fmt.Println("启动web服务失败！")
	}
}
