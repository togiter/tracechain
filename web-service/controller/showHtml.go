package controller

import (
	"fmt"
	"net/http"
)

func (app *Application) ShowIssueHtml(w http.ResponseWriter, r *http.Request){
	fmt.Println("Issue html")
	renderTemplate(w, r, "issue", nil)
}