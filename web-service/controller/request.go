package controller

import (
	"fmt"
	"net/http"
	// "github.com/tracechain/fabric-service/fabricSetup"
)

func (app *Application) IssueProduct(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	fmt.Println("requesting Url...", r.URL)
	fmt.Println("请求参数:")
	name := r.PostFormValue("name")
	number := r.PostFormValue("number")
	millPrice := r.PostFormValue("millPrice")
	price := r.PostFormValue("price")
	color := r.PostFormValue("color")
	owner := r.PostFormValue("owner")
	productor := r.PostFormValue("productor")
	fmt.Printf("name:%s\nnumber:%s\nmillPrice:%s\nprice:%s\ncolor:%s\nowner:%s\nproductor:%s\n", name,number, millPrice, price, color, owner, productor)
	result, err := app.Fabric.IssueProduct(name, number, millPrice, price, color, owner, productor)
	if err != nil {
		fmt.Println("error:", err)
		renderTemplate(w, r, "mainlayout", err)
		return
	}
	fmt.Println("Issue TX", result)
	renderTemplate(w, r, "mainlayout", result)
}

func (app *Application) TransferProduct(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	fmt.Println("transfering product...")
	number := r.PostFormValue("number")
	owner := r.PostFormValue("owner")
	price := r.PostFormValue("price")
	fmt.Printf("owner:%s\nnumber:%s\nprice:%s\n", owner, number, price)
	result, err := app.Fabric.TransferProduct(owner, number, price)
	if err != nil {
		fmt.Println("error:", err)
		renderTemplate(w, r, "mainlayout", err)
		return
	}
	fmt.Println("transfer TX:", result)
	renderTemplate(w, r, "mainlayout", result)

}

func (app *Application) AlterProductPrice(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	fmt.Println("altering product price...")
	number := r.PostFormValue("number")
	owner := r.PostFormValue("owner")
	price := r.PostFormValue("price")
	fmt.Printf("owner:%s\nnumber:%s\nprice:%s\n", owner, number, price)
	result, err := app.Fabric.AlterProductPrice(owner, number, price)
	if err != nil {
		fmt.Println("error:", err)
		renderTemplate(w, r, "mainlayout", err)
		return
	}
	fmt.Println("alter TX:", result)
	renderTemplate(w, r, "mainlayout", result)
}

func (app *Application) QueryProductNo(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	number := r.Form["number"][0]
	fmt.Println(" query product with number:", number)
	result, err := app.Fabric.QueryProductNo(number)
	if err != nil {
		fmt.Println("error:", err)
		renderTemplate(w, r, "mainlayout", err)
		return
	}
	fmt.Println("query result:", result)
	renderTemplate(w, r, "mainlayout", result)
}

func (app *Application) QueryProductRange(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	startkey := r.Form["startKey"][0]
	endKey := r.Form["endKey"][0]
	fmt.Printf("query range condition:%s~%s", startkey, endKey)
	result, err := app.Fabric.QueryProductRange(startkey, endKey)
	if err != nil {
		fmt.Println("error:", err)
		renderTemplate(w, r, "mainlayout", err)
		return
	}
	fmt.Println("query range result:", result)
	renderTemplate(w, r, "mainlayout", result)
}
