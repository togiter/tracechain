package controller

import (
	"fmt"
	"net/http"
	"github.com/tracechain/fabric-service/product"
)

func (app *Application)IssueProduct(w http.ResponseWriter,r *http.Request){
	r.ParseForm()
	fmt.Println("requesting Url...",r.URL)
	fmt.Println("请求参数:")
	name := r.Form["name"]
	number := r.Form["number"]
	millPrice := r.Form["millPrice"]
	price := r.Form["price"]
	color := r.Form["color"]
	owner := r.Form["owner"]
	productor := r.Form["productor"]
	fmt.Printf("name:%s\nnumber:%s\nmillPrice:%s\nprice:%s\ncolor:%s\nowner:%s\nproductor:%s\n",name,millPrice,price,color,owner,productor)
	result,err := productservice.Issue(app.Fabric,name,number,millPrice,price,color,owner,productor)
	if err != nil {
		fmt.Println("error:",err)
		renderTemplate(w,r,"mainlayout",err)
		return
	}
	fmt.Println("Issue TX",result)
	renderTemplate(w,r,"mainlayout",result)
}

func (app *Application)TransferProduct(w http.ResponseWriter,r *http.ReadRequest){
	r.ParseForm()
	fmt.Println("transfering product...")
	number := r.Form["number"]
	owner := r.Form["owner"]
	price := r.Form["price"]
	fmt.Printf("owner:%s\nnumber:%s\nprice:%s\n",owner,number,price)
	result,err := productservice.TransferProduct(app.Fabric,owner,number,price)
	if err != nil {
		fmt.Println("error:",err)
		renderTemplate(w,r,"mainlayout",err)
		return
	}
	fmt.Println("transfer TX:",result)
	renderTemplate(w,r,"mainlayout",result)

}

func (app *Application)AlterPrice(w http.ResponseWriter,r *http.Request){
	r.ParseForm()
	fmt.Println("altering product price...")
	number := r.Form["number"]
	owner := r.Form["owner"]
	price := r.Form["price"]
	fmt.Printf("owner:%s\nnumber:%s\nprice:%s\n",owner,number,price)
	result,err := productservice.AlterProductPrice(app.Fabric,owner,number,price)
	if err != nil {
		fmt.Println("error:",err)
		renderTemplate(w,r,"mainlayout",err)
		return
	}
	fmt.Println("alter TX:",result)
	renderTemplate(w,r,"mainlayout",result)
}

func (app *Application)QueryProductNo(w http.ResponseWriter,r *http.Request){
	r.ParseForm()
	number := r.Form["number"]
	fmt.Println(" query product with number:",number)
	result,err := productservice.QueryProductNo(app.Fabric,number)
	if err != nil {
		fmt.Println("error:",err)
		renderTemplate(w,r,"mainlayout",err)
		return
	}
	fmt.Println("query result:",result)
	renderTemplate(w,r,"mainlayout",result)
}

func (app *Application)QueryProductRange(w http.ResponseWriter,r *http.Request){
	r.ParseForm()
	startkey := r.Form["startKey"]
	endKey := r.Form["endKey"]
	fmt.Printf("query range condition:%s~%s",startkey,endKey)
	result,err := productservice.QueryProductRange(app.Fabric,startkey,endKey)
	if err != nil {
		fmt.Println("error:",err)
		renderTemplate(w,r,"mainlayout",err)
		return
	}
	fmt.Println("query range result:",result)
	renderTemplate(w,r,"mainlayout",result)
}